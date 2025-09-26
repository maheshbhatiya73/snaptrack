package backups

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"snaptrack/db"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type BackupService struct {
    clients   map[*websocket.Conn]bool
    clientsMu sync.RWMutex
}


func timePtr(t time.Time) *time.Time { return &t }


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
                bs.updateProgress(progress, progress.Progress, "failed", fmt.Sprintf("Local backup failed: %v", err))
                // Persist backup failed status
                backup.Status = "failed"
                backup.CompletedAt = timePtr(time.Now())
                backup.DurationSec = int64(time.Since(startTime).Seconds())
                db.DB.Save(&backup)
                db.DB.Create(&db.Log{Level: "error", Message: fmt.Sprintf("Backup %s failed: %v", backup.Name, err)})
                return
            }
			backup.SizeBytes = totalSize
			backup.Checksum = &checksum

		} else if server.Type == "remote" {
            err := bs.runRemoteBackup(backup, server, progress)
            if err != nil {
                bs.updateProgress(progress, progress.Progress, "failed", fmt.Sprintf("Remote backup failed: %v", err))
                // Persist backup failed status
                backup.Status = "failed"
                backup.CompletedAt = timePtr(time.Now())
                backup.DurationSec = int64(time.Since(startTime).Seconds())
                db.DB.Save(&backup)
                db.DB.Create(&db.Log{Level: "error", Message: fmt.Sprintf("Backup %s failed: %v", backup.Name, err)})
                return
            }
		}
	}

    now := time.Now()
    if progress.Status != "failed" { // Do not override failed state
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

func (bs *BackupService) updateProgress(progress *db.BackupProgress, percentage int, status, message string) {
	progress.Progress = percentage
	progress.Status = status
	progress.Message = message
	progress.UpdatedAt = time.Now()
	db.DB.Save(progress)
	bs.BroadcastProgress(progress)
}