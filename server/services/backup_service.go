package services

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"snaptrack/db"
	"strconv"
	"strings"
	"sync"
	"time"

	"archive/tar"
	"archive/zip"

	"github.com/gofiber/websocket/v2"
	"golang.org/x/crypto/ssh"
)

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

type BackupService struct {
	clients   map[*websocket.Conn]bool
	clientsMu sync.RWMutex
}

type countingReader struct {
	r      io.Reader
	onRead func(int64)
}

func (cr countingReader) Read(p []byte) (int, error) {
	n, err := cr.r.Read(p)
	if n > 0 {
		cr.onRead(int64(n))
	}
	return n, err
}

func NewBackupService() *BackupService {
	return &BackupService{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (bs *BackupService) AddClient(c *websocket.Conn) {
	bs.clientsMu.Lock()
	bs.clients[c] = true
	bs.clientsMu.Unlock()
}

func (bs *BackupService) RemoveClient(c *websocket.Conn) {
	bs.clientsMu.Lock()
	delete(bs.clients, c)
	bs.clientsMu.Unlock()
}

func (bs *BackupService) BroadcastProgress(progress *db.BackupProgress) {
	fmt.Printf("[BACKUP PROGRESS] BackupID: %d, Status: %s, Progress: %d%%, Message: %s\n", progress.BackupID, progress.Status, progress.Progress, progress.Message)

	var backup db.Backup
	if err := db.DB.First(&backup, progress.BackupID).Error; err != nil {
		fmt.Printf("[BACKUP PROGRESS ERROR] Failed to load backup %d: %v\n", progress.BackupID, err)
		return
	}
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
			continue
		}
	}
}

func (bs *BackupService) ExecuteBackupAsync(backup db.Backup) {
	startTime := time.Now()
	progress := &db.BackupProgress{
		BackupID: backup.ID,
		Status:   "running",
		Progress: 0,
		Message:  "Starting backup...",
	}
	db.DB.Create(progress)
	bs.BroadcastProgress(progress)

	var serverIDs []uint
	if len(backup.ServerIDs) > 0 {
		if err := json.Unmarshal(backup.ServerIDs, &serverIDs); err != nil {
			bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Invalid server_ids format: %v", err))
			return
		}
	}

	hasRemote := false
	for _, serverID := range serverIDs {
		var server db.Server
		if err := db.DB.First(&server, serverID).Error; err != nil {
			bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Server %d not found: %v", serverID, err))
			return
		}

		if server.Type == "remote" {
			hasRemote = true
			if err := validateRemoteServer(server.Host, server.SSHUser, server.SSHKeyPath, server.SSHPort, server.TransferType); err != nil {
				bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Remote server validation failed: %v", err))
				return
			}
		}
	}

	if hasRemote {
		if _, err := exec.LookPath("rsync"); err != nil {
			bs.updateProgress(progress, 0, "failed", "rsync command not found in PATH")
			return
		}
	}

	// Validate source path
	if _, err := os.Stat(backup.Source); os.IsNotExist(err) {
		bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Source path does not exist: %s", backup.Source))
		return
	}

	// Handle each server separately
	for _, serverID := range serverIDs {
		var server db.Server
		db.DB.First(&server, serverID)

		if server.Type == "local" {
			totalSize, checksum, err := bs.runLocalBackup(backup, progress)
			if err != nil {
				bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Local backup failed: %v", err))
				return
			}
			backup.SizeBytes = totalSize
			backup.Checksum = &checksum

		} else if server.Type == "remote" {
			err := bs.runRemoteBackup(backup, server, progress)
			if err != nil {
				bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Remote backup failed: %v", err))
				return
			}
		}
	}

	now := time.Now()
	backup.Status = "completed"
	backup.CompletedAt = &now
	backup.DurationSec = int64(now.Sub(startTime).Seconds())
	db.DB.Save(&backup)
	db.DB.Create(&db.Log{
		Level:   "info",
		Message: fmt.Sprintf("Backup %s completed successfully", backup.Name),
	})
	bs.updateProgress(progress, 100, "completed", "Backup completed successfully")
}

func validateRemoteServer(host string, user *string, keyPath *string, port *int, transferType *string) error {
	if host == "" || user == nil || keyPath == nil || port == nil {
		return fmt.Errorf("missing SSH credentials or host")
	}

	key, err := os.ReadFile(*keyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse SSH key: %v", err)
	}

	config := &ssh.ClientConfig{
		User:            *user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, *port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("cannot connect to remote server: %v", err)
	}
	client.Close()
	return nil
}

func computeTotalSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func (bs *BackupService) runLocalBackup(backup db.Backup, progress *db.BackupProgress) (int64, string, error) {
	switch backup.FileType {
	case "tar":
		return bs.createTarBackupWithProgress(backup.Source, backup.Destination, progress)
	case "zip":
		return bs.createZipBackupWithProgress(backup.Source, backup.Destination, progress)
	case "raw":
		return bs.createRawBackupWithProgress(backup.Source, backup.Destination, progress)
	default:
		return 0, "", fmt.Errorf("unsupported file type: %s", backup.FileType)
	}
}

func (bs *BackupService) runRemoteBackup(backup db.Backup, server db.Server, progress *db.BackupProgress) error {
	if backup.FileType == "raw" {
		source := backup.Source + "/."
		dest := backup.Destination
		total, err := bs.getRsyncTotalSize(source, server, dest)
		if err != nil {
			return err
		}
		t := total
		progress.TotalBytes = &t
		progress.BytesProcessed = 0
		err = bs.runRsyncBackup(source, server, dest, progress, 0)
		if err != nil {
			return err
		}
		checksum, err := bs.calculateDirectoryChecksum(backup.Source)
		if err != nil {
			return err
		}
		backup.SizeBytes = total
		backup.Checksum = &checksum
		db.DB.Save(&backup)
		return nil
	} else {
		tempDir, err := os.MkdirTemp("", "backup-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)

		totalUncomp, err := computeTotalSize(backup.Source)
		if err != nil {
			return err
		}
		t := totalUncomp
		progress.TotalBytes = &t
		progress.BytesProcessed = 0
		bs.updateProgress(progress, 0, "running", "Archiving locally...")
		var archSize int64
		var checksum string
		var archFile string
		if backup.FileType == "tar" {
			archSize, checksum, err = bs.createTarBackupWithProgress(backup.Source, tempDir, progress)
			archFile = filepath.Join(tempDir, "backup.tar")
		} else if backup.FileType == "zip" {
			archSize, checksum, err = bs.createZipBackupWithProgress(backup.Source, tempDir, progress)
			archFile = filepath.Join(tempDir, "backup.zip")
		} else {
			return fmt.Errorf("unsupported file type: %s", backup.FileType)
		}
		if err != nil {
			return err
		}
		archProcessed := progress.BytesProcessed
		overallTotal := archProcessed + archSize
		ot := overallTotal
		progress.TotalBytes = &ot
		progress.Progress = int(archProcessed * 100 / overallTotal)
		progress.Message = "Transferring to remote..."
		progress.BytesProcessed = archProcessed
		db.DB.Save(progress)
		bs.BroadcastProgress(progress)
		err = bs.runRsyncBackup(archFile, server, backup.Destination, progress, archProcessed)
		if err != nil {
			return err
		}
		backup.SizeBytes = archSize
		backup.Checksum = &checksum
		db.DB.Save(&backup)
		return nil
	}
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
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}
	tarFilePath := filepath.Join(destinationDir, "backup.tar")
	f, err := os.Create(tarFilePath)
	if err != nil {
		return 0, "", err
	}
	gzw := gzip.NewWriter(f)
	tw := tar.NewWriter(gzw)
	start := time.Now()
	lastUpdate := time.Now()
	delta := int64(0)
	onRead := func(n int64) {
		progress.BytesProcessed += n
		delta += n
		now := time.Now()
		if delta >= 1<<20 || now.Sub(lastUpdate) >= time.Second {
			elapsed := now.Sub(start).Seconds()
			var speed int64
			if elapsed > 0 {
				speed = int64(float64(progress.BytesProcessed) / elapsed)
			}
			var eta int64
			if speed > 0 {
				eta = (*progress.TotalBytes - progress.BytesProcessed) / speed
			}
			progress.SpeedBPS = &speed
			progress.ETASeconds = &eta
			progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
			progress.UpdatedAt = now
			db.DB.Save(progress)
			bs.BroadcastProgress(progress)
			lastUpdate = now
			delta = 0
		}
	}
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		header.Name = rel
		if err = tw.WriteHeader(header); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		curr := filepath.Base(path)
		progress.CurrentFile = &curr
		bs.BroadcastProgress(progress)
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, countingReader{r: in, onRead: onRead})
		in.Close()
		return err
	})
	if err != nil {
		tw.Close()
		gzw.Close()
		f.Close()
		return 0, "", err
	}
	if err = tw.Close(); err != nil {
		gzw.Close()
		f.Close()
		return 0, "", err
	}
	if err = gzw.Close(); err != nil {
		f.Close()
		return 0, "", err
	}
	f.Close()
	now := time.Now()
	elapsed := now.Sub(start).Seconds()
	var speed int64
	if elapsed > 0 {
		speed = int64(float64(progress.BytesProcessed) / elapsed)
	}
	eta0 := int64(0)
	progress.SpeedBPS = &speed
	progress.ETASeconds = &eta0
	progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
	progress.UpdatedAt = now
	db.DB.Save(progress)
	bs.BroadcastProgress(progress)
	bs.updateProgress(progress, progress.Progress, "running", "Calculating checksum...")
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
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}
	zipFilePath := filepath.Join(destinationDir, "backup.zip")
	f, err := os.Create(zipFilePath)
	if err != nil {
		return 0, "", err
	}
	zw := zip.NewWriter(f)
	start := time.Now()
	lastUpdate := time.Now()
	delta := int64(0)
	onRead := func(n int64) {
		progress.BytesProcessed += n
		delta += n
		now := time.Now()
		if delta >= 1<<20 || now.Sub(lastUpdate) >= time.Second {
			elapsed := now.Sub(start).Seconds()
			var speed int64
			if elapsed > 0 {
				speed = int64(float64(progress.BytesProcessed) / elapsed)
			}
			var eta int64
			if speed > 0 {
				eta = (*progress.TotalBytes - progress.BytesProcessed) / speed
			}
			progress.SpeedBPS = &speed
			progress.ETASeconds = &eta
			progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
			progress.UpdatedAt = now
			db.DB.Save(progress)
			bs.BroadcastProgress(progress)
			lastUpdate = now
			delta = 0
		}
	}
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		header := &zip.FileHeader{
			Name:     rel,
			Method:   zip.Deflate,
			Modified: info.ModTime(),
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		curr := filepath.Base(path)
		progress.CurrentFile = &curr
		bs.BroadcastProgress(progress)
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, countingReader{r: in, onRead: onRead})
		in.Close()
		return err
	})
	if err != nil {
		zw.Close()
		f.Close()
		return 0, "", err
	}
	if err = zw.Close(); err != nil {
		f.Close()
		return 0, "", err
	}
	f.Close()
	now := time.Now()
	elapsed := now.Sub(start).Seconds()
	var speed int64
	if elapsed > 0 {
		speed = int64(float64(progress.BytesProcessed) / elapsed)
	}
	eta0 := int64(0)
	progress.SpeedBPS = &speed
	progress.ETASeconds = &eta0
	progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
	progress.UpdatedAt = now
	db.DB.Save(progress)
	bs.BroadcastProgress(progress)
	bs.updateProgress(progress, progress.Progress, "running", "Calculating checksum...")
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
	if err := os.MkdirAll(destination, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}
	start := time.Now()
	lastUpdate := time.Now()
	delta := int64(0)
	onRead := func(n int64) {
		progress.BytesProcessed += n
		delta += n
		now := time.Now()
		if delta >= 1<<20 || now.Sub(lastUpdate) >= time.Second {
			elapsed := now.Sub(start).Seconds()
			var speed int64
			if elapsed > 0 {
				speed = int64(float64(progress.BytesProcessed) / elapsed)
			}
			var eta int64
			if speed > 0 {
				eta = (*progress.TotalBytes - progress.BytesProcessed) / speed
			}
			progress.SpeedBPS = &speed
			progress.ETASeconds = &eta
			progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
			progress.UpdatedAt = now
			db.DB.Save(progress)
			bs.BroadcastProgress(progress)
			lastUpdate = now
			delta = 0
		}
	}
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destination, rel)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		curr := filepath.Base(path)
		progress.CurrentFile = &curr
		bs.BroadcastProgress(progress)
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		out, err := os.Create(destPath)
		if err != nil {
			in.Close()
			return err
		}
		_, err = io.Copy(out, countingReader{r: in, onRead: onRead})
		in.Close()
		out.Close()
		return err
	})
	if err != nil {
		return 0, "", err
	}
	now := time.Now()
	elapsed := now.Sub(start).Seconds()
	var speed int64
	if elapsed > 0 {
		speed = int64(float64(progress.BytesProcessed) / elapsed)
	}
	eta0 := int64(0)
	progress.SpeedBPS = &speed
	progress.ETASeconds = &eta0
	progress.Progress = 100
	progress.UpdatedAt = now
	db.DB.Save(progress)
	bs.BroadcastProgress(progress)
	bs.updateProgress(progress, 100, "running", "Calculating checksum...")
	checksum, err := bs.calculateDirectoryChecksum(destination)
	if err != nil {
		return 0, "", fmt.Errorf("failed to calculate checksum: %v", err)
	}
	return *progress.TotalBytes, checksum, nil
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

func (bs *BackupService) GetBackupProgress(backupID uint) (*db.BackupProgress, error) {
	var progress db.BackupProgress
	err := db.DB.Where("backup_id = ?", backupID).Order("updated_at DESC").First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (bs *BackupService) GetAllRunningBackups() ([]db.BackupProgress, error) {
	var progresses []db.BackupProgress
	err := db.DB.Where("status IN ?", []string{"running", "failed", "completed"}).Preload("Backup").Find(&progresses).Error
	return progresses, err
}

func (bs *BackupService) getRsyncTotalSize(source string, server db.Server, dest string) (int64, error) {
	args := []string{
		"-az",
		"--dry-run",
		"--stats",
		source,
		fmt.Sprintf("%s@%s:%s", *server.SSHUser, server.Host, dest),
	}
	cmd := exec.Command("rsync", args...)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("rsync dry-run failed: %v", err)
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Total file size:") {
			sizeStr := strings.TrimSpace(strings.TrimPrefix(line, "Total file size:"))
			sizeStr = strings.TrimSuffix(sizeStr, " bytes")
			sizeStr = strings.ReplaceAll(sizeStr, ",", "")
			return strconv.ParseInt(sizeStr, 10, 64)
		}
	}
	return 0, fmt.Errorf("could not find total file size in rsync output")
}

func (bs *BackupService) runRsyncBackup(source string, server db.Server, dest string, progress *db.BackupProgress, offsetBytes int64) error {
	curr := "overall"
	progress.CurrentFile = &curr
	bs.updateProgress(progress, progress.Progress, "running", fmt.Sprintf("Starting transfer to %s", server.Host))
	args := []string{
		"-az",
		"--info=progress2",
		source,
		fmt.Sprintf("%s@%s:%s", *server.SSHUser, server.Host, dest),
	}
	cmd := exec.Command("rsync", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start rsync: %v", err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		processed, percent, speed, eta := parseRsyncProgress(line)
		if percent >= 0 {
			progress.BytesProcessed = offsetBytes + processed
			progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
			progress.SpeedBPS = &speed
			progress.ETASeconds = &eta
			progress.UpdatedAt = time.Now()
			db.DB.Save(progress)
			bs.BroadcastProgress(progress)
		}
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("rsync failed: %v", err)
	}
	bs.updateProgress(progress, 100, "completed", "Transfer completed successfully")
	return nil
}

func parseRsyncProgress(line string) (int64, int, int64, int64) {
	line = strings.TrimSpace(line)
	fields := strings.Fields(line)
	if len(fields) < 4 || !strings.HasSuffix(fields[1], "%") {
		return 0, -1, 0, 0
	}
	bytesStr := strings.ReplaceAll(fields[0], ",", "")
	processed, err := strconv.ParseInt(bytesStr, 10, 64)
	if err != nil {
		return 0, -1, 0, 0
	}
	percentStr := strings.TrimSuffix(fields[1], "%")
	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		return 0, -1, 0, 0
	}
	speedStr := fields[2]
	unit := speedStr[len(speedStr)-3:]
	speedNum, err := strconv.ParseFloat(speedStr[:len(speedStr)-3], 64)
	if err != nil {
		return 0, -1, 0, 0
	}
	var multi float64 = 1
	switch unit {
	case "B/s":
		multi = 1
	case "KB/s":
		multi = 1024
	case "MB/s":
		multi = 1024 * 1024
	case "GB/s":
		multi = 1024 * 1024 * 1024
	}
	speed := int64(speedNum * multi)
	etaStr := fields[3]
	etaParts := strings.Split(etaStr, ":")
	if len(etaParts) != 3 {
		return processed, percent, speed, 0
	}
	h, _ := strconv.Atoi(etaParts[0])
	m, _ := strconv.Atoi(etaParts[1])
	s, _ := strconv.Atoi(etaParts[2])
	eta := int64(h*3600 + m*60 + s)
	return processed, percent, speed, eta
}