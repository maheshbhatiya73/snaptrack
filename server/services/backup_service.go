package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"snaptrack/db"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

// BackupProgressResponse represents the progress data sent to WebSocket clients
type BackupProgressResponse struct {
	ID             uint        `json:"id"`
	BackupID       uint        `json:"backup_id"`
	Backup         db.Backup   `json:"backup"`
	Servers        []db.Server `json:"servers"`
	Status         string      `json:"status"`
	Progress       int         `json:"progress"`
	Message        string      `json:"message"`
	CurrentFile    *string     `json:"current_file"`
	BytesProcessed int64       `json:"bytes_processed"`
	TotalBytes     *int64      `json:"total_bytes"`
	SpeedBPS       *int64      `json:"speed_bps"`
	ETASeconds     *int64      `json:"eta_seconds"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// BackupService manages backup operations and progress tracking
type BackupService struct {
	clients   map[*websocket.Conn]bool
	clientsMu sync.RWMutex
}

// NewBackupService creates a new backup service instance
func NewBackupService() *BackupService {
	return &BackupService{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient adds a WebSocket client for real-time updates
func (bs *BackupService) AddClient(c *websocket.Conn) {
	bs.clientsMu.Lock()
	bs.clients[c] = true
	bs.clientsMu.Unlock()
}

// RemoveClient removes a WebSocket client
func (bs *BackupService) RemoveClient(c *websocket.Conn) {
	bs.clientsMu.Lock()
	delete(bs.clients, c)
	bs.clientsMu.Unlock()
}

// BroadcastProgress sends progress updates to all connected clients
func (bs *BackupService) BroadcastProgress(progress *db.BackupProgress) {
	fmt.Printf("[BACKUP PROGRESS] BackupID: %d, Status: %s, Progress: %d%%, Message: %s\n", progress.BackupID, progress.Status, progress.Progress, progress.Message)

	// Load the backup data
	var backup db.Backup
	if err := db.DB.First(&backup, progress.BackupID).Error; err != nil {
		fmt.Printf("[BACKUP PROGRESS ERROR] Failed to load backup %d: %v\n", progress.BackupID, err)
		return
	}

	// Load server data
	var serverIDs []uint
	if len(backup.ServerIDs) > 0 {
		if err := json.Unmarshal(backup.ServerIDs, &serverIDs); err != nil {
			fmt.Printf("[BACKUP PROGRESS ERROR] Failed to unmarshal server_ids: %v\n", err)
		}
	}
	var servers []db.Server
	if len(serverIDs) > 0 {
		db.DB.Find(&servers, serverIDs)
	}

	response := BackupProgressResponse{
		ID:             progress.ID,
		BackupID:       progress.BackupID,
		Backup:         backup,
		Servers:        servers,
		Status:         progress.Status,
		Progress:       progress.Progress,
		Message:        progress.Message,
		CurrentFile:    progress.CurrentFile,
		BytesProcessed: progress.BytesProcessed,
		TotalBytes:     progress.TotalBytes,
		SpeedBPS:       progress.SpeedBPS,
		ETASeconds:     progress.ETASeconds,
		CreatedAt:      progress.CreatedAt,
		UpdatedAt:      progress.UpdatedAt,
	}

	bs.clientsMu.RLock()
	defer bs.clientsMu.RUnlock()

	for client := range bs.clients {
		if err := client.WriteJSON(response); err != nil {
			// Client disconnected, will be cleaned up by the connection handler
			continue
		}
	}
}

// ExecuteBackupAsync executes a backup with progress tracking
func (bs *BackupService) ExecuteBackupAsync(backup db.Backup) {
	startTime := time.Now()
	var totalSize int64
	var checksum string

	// Create initial progress entry
	progress := &db.BackupProgress{
		BackupID: backup.ID,
		Status:   "running",
		Progress: 0,
		Message:  "Starting backup...",
	}
	db.DB.Create(progress)
	bs.BroadcastProgress(progress)

	// Validate servers
	var serverIDs []uint
	if len(backup.ServerIDs) > 0 {
		if err := json.Unmarshal(backup.ServerIDs, &serverIDs); err != nil {
			bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Invalid server_ids format: %v", err))
			return
		}
	}

	// For now, focus on local backups
	for _, serverID := range serverIDs {
		var server db.Server
		if err := db.DB.First(&server, serverID).Error; err != nil {
			bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Server %d not found: %v", serverID, err))
			return
		}

		if server.Type != "local" {
			bs.updateProgress(progress, 0, "failed", "Remote servers not supported yet")
			return
		}
	}

	// Check if required commands are available
	switch backup.FileType {
	case "tar":
		if _, err := exec.LookPath("tar"); err != nil {
			bs.updateProgress(progress, 0, "failed", "tar command not found in PATH")
			return
		}
	case "zip":
		if _, err := exec.LookPath("zip"); err != nil {
			bs.updateProgress(progress, 0, "failed", "zip command not found in PATH")
			return
		}
	}

	bs.updateProgress(progress, 10, "running", "Validating paths...")

	// Validate source and destination paths
	// Check if source exists
	if _, err := os.Stat(backup.Source); os.IsNotExist(err) {
		bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Source path does not exist: %s", backup.Source))
		return
	}

	// Check if we can write to destination directory
	destDir := backup.Destination
	if err := os.MkdirAll(destDir, 0755); err != nil {
		bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Cannot create destination directory: %v", err))
		return
	}

	// Check if destination is writable
	testFile := filepath.Join(destDir, ".backup_test")
	if f, err := os.Create(testFile); err != nil {
		bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Destination is not writable: %v", err))
		return
	} else {
		f.Close()
		os.Remove(testFile)
	}

	// Perform backup operation based on file type
	var err error

	switch backup.FileType {
	case "tar":
		totalSize, checksum, err = bs.createTarBackupWithProgress(backup.Source, backup.Destination, progress)
	case "zip":
		totalSize, checksum, err = bs.createZipBackupWithProgress(backup.Source, backup.Destination, progress)
	case "raw":
		totalSize, checksum, err = bs.createRawBackupWithProgress(backup.Source, backup.Destination, progress)
	default:
		err = fmt.Errorf("unsupported file type: %s", backup.FileType)
	}

	if err != nil {
		bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Backup failed: %v", err))
		// Update backup status
		backup.Status = "failed"
		backup.CompletedAt = &startTime
		backup.DurationSec = int64(time.Since(startTime).Seconds())
		db.DB.Save(&backup)
		return
	}

	// Update backup record for success
	backup.Status = "completed"
	backup.CompletedAt = &startTime
	backup.DurationSec = int64(time.Since(startTime).Seconds())
	backup.SizeBytes = totalSize
	backup.Checksum = &checksum
	db.DB.Save(&backup)

	// Log successful backup
	db.DB.Create(&db.Log{
		Level:   "info",
		Message: fmt.Sprintf("Backup %s completed successfully", backup.Name),
	})

	bs.updateProgress(progress, 100, "completed", "Backup completed successfully")
}

func (bs *BackupService) updateProgress(progress *db.BackupProgress, percentage int, status, message string) {
	progress.Progress = percentage
	progress.Status = status
	progress.Message = message
	progress.UpdatedAt = time.Now()
	db.DB.Save(progress)
	bs.BroadcastProgress(progress)
}

func (bs *BackupService) createTarBackupWithProgress(source, destinationDir string, progress *db.BackupProgress) (int64, string, error) {
	bs.updateProgress(progress, 20, "running", "Creating tar archive...")

	// Ensure destination directory exists
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Create the tar file path
	tarFilePath := filepath.Join(destinationDir, "backup.tar")

	// Create tar command
	cmd := exec.Command("tar", "-czf", tarFilePath, "-C", filepath.Dir(source), filepath.Base(source))
	if err := cmd.Run(); err != nil {
		return 0, "", fmt.Errorf("tar command failed: %v", err)
	}

	bs.updateProgress(progress, 80, "running", "Calculating checksum...")

	// Get file size and checksum
	info, err := os.Stat(tarFilePath)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get file info: %v", err)
	}
	size := info.Size()

	checksum, err := bs.calculateFileChecksum(tarFilePath)
	if err != nil {
		return 0, "", fmt.Errorf("failed to calculate checksum: %v", err)
	}

	return size, checksum, nil
}

func (bs *BackupService) createZipBackupWithProgress(source, destinationDir string, progress *db.BackupProgress) (int64, string, error) {
	bs.updateProgress(progress, 20, "running", "Creating zip archive...")

	// Ensure destination directory exists
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Create the zip file path
	zipFilePath := filepath.Join(destinationDir, "backup.zip")

	// Create zip command
	cmd := exec.Command("zip", "-r", zipFilePath, filepath.Base(source))
	cmd.Dir = filepath.Dir(source)
	if err := cmd.Run(); err != nil {
		return 0, "", fmt.Errorf("zip command failed: %v", err)
	}

	bs.updateProgress(progress, 80, "running", "Calculating checksum...")

	// Get file size and checksum
	info, err := os.Stat(zipFilePath)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get file info: %v", err)
	}
	size := info.Size()

	checksum, err := bs.calculateFileChecksum(zipFilePath)
	if err != nil {
		return 0, "", fmt.Errorf("failed to calculate checksum: %v", err)
	}

	return size, checksum, nil
}

func (bs *BackupService) createRawBackupWithProgress(source, destination string, progress *db.BackupProgress) (int64, string, error) {
	bs.updateProgress(progress, 20, "running", "Copying files...")

	// Ensure destination directory exists
	if err := os.MkdirAll(destination, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Copy files recursively
	cmd := exec.Command("cp", "-r", source+"/.", destination+"/")
	if err := cmd.Run(); err != nil {
		return 0, "", fmt.Errorf("cp command failed: %v", err)
	}

	bs.updateProgress(progress, 80, "running", "Calculating size and checksum...")

	// Calculate total size
	var totalSize int64
	err := filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("failed to calculate size: %v", err)
	}

	// Calculate checksum
	checksum, err := bs.calculateDirectoryChecksum(destination)
	if err != nil {
		return 0, "", fmt.Errorf("failed to calculate checksum: %v", err)
	}

	return totalSize, checksum, nil
}

func (bs *BackupService) calculateFileChecksum(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func (bs *BackupService) calculateDirectoryChecksum(dirPath string) (string, error) {
	hash := sha256.New()

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			hash.Write(data)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetBackupProgress retrieves the latest progress for a backup
func (bs *BackupService) GetBackupProgress(backupID uint) (*db.BackupProgress, error) {
	var progress db.BackupProgress
	err := db.DB.Where("backup_id = ?", backupID).Order("updated_at DESC").First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// GetAllRunningBackups returns all currently running, failed, and recently completed backups with their progress
func (bs *BackupService) GetAllRunningBackups() ([]db.BackupProgress, error) {
	var progresses []db.BackupProgress
	err := db.DB.Where("status IN ?", []string{"running", "failed", "completed"}).Preload("Backup").Find(&progresses).Error
	return progresses, err
}