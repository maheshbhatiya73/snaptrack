package routes

import (
	"encoding/json"
	"snaptrack/auth"
	"snaptrack/db"
	"snaptrack/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// BackupProgressResponse represents the progress data sent to clients
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

type BackupResponse struct {
	ID           uint        `json:"id"`
	Name         string      `json:"name"`
	ServerIDs    []uint      `json:"server_ids"`
	Servers      []db.Server `json:"servers"` 
	Source       string      `json:"source"`
	Destination  string      `json:"destination"`
	FileType     string      `json:"file_type"`
	Type         string      `json:"type"`
	ScheduleType string      `json:"schedule_type"`
	Status       string      `json:"status"`
	SizeBytes    int64       `json:"size_bytes"`
	Checksum     *string     `json:"checksum"`
	StartedAt    *time.Time  `json:"started_at"`
	CompletedAt  *time.Time  `json:"completed_at"`
	DurationSec  int64       `json:"duration_sec"`
	ExecutedBy   string      `json:"executed_by"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

var backupService *services.BackupService

func RegisterBackupRoutes(app *fiber.App) {
	backupService = services.NewBackupService()

    // On service start, clean up any stale running states from previous instance
    cleanupStaleRunning()

	api := app.Group("/api/backups", auth.RequireJWT())

	api.Get("/", listBackups)
	api.Post("/", createBackup)
	api.Get("/processes/running", getRunningBackups)
	api.Delete("/processes", deleteAllProcesses)
	api.Delete("/processes/:id", deleteProcess)
	api.Get("/:id", getBackup)
	api.Put("/:id", updateBackup)
	api.Delete("/:id", deleteBackup)
	api.Post("/:id/execute", executeBackup)
	api.Get("/:id/progress", getBackupProgress)
}

// cleanupStaleRunning marks any "running" progress/backup records as interrupted
// so the system can be re-executed after a service restart or crash.
func cleanupStaleRunning() {
    // Mark running progresses as failed with a restart message
    var progresses []db.BackupProgress
    if err := db.DB.Where("status = ?", "running").Find(&progresses).Error; err == nil {
        for _, p := range progresses {
            p.Status = "failed"
            p.Message = "Interrupted due to service restart"
            p.UpdatedAt = time.Now()
            db.DB.Save(&p)

            // Reset corresponding backup to pending so it can be started again
            var b db.Backup
            if err := db.DB.First(&b, p.BackupID).Error; err == nil {
                b.Status = "pending"
                b.UpdatedAt = time.Now()
                db.DB.Save(&b)
            }
        }
    }

    // Also reset any backups stuck in running without a progress row
    db.DB.Model(&db.Backup{}).Where("status = ?", "running").Updates(map[string]any{
        "status":     "pending",
        "updated_at": time.Now(),
    })
}

// WebSocket route for real-time updates
func RegisterWebSocketRoutes(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/backups", websocket.New(func(c *websocket.Conn) {
		// Add client to service
		backupService.AddClient(c)
		defer backupService.RemoveClient(c)

		// Keep connection alive
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))
}

func listBackups(c *fiber.Ctx) error {
	var backups []db.Backup
	db.DB.Find(&backups)

	var resp []BackupResponse
	for _, b := range backups {
		var serverIDs []uint
		if len(b.ServerIDs) > 0 {
			if err := json.Unmarshal(b.ServerIDs, &serverIDs); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Invalid server_ids format"})
			}
		}

		var servers []db.Server
		if len(serverIDs) > 0 {
			db.DB.Find(&servers, serverIDs)
		}

		resp = append(resp, BackupResponse{
			ID:           b.ID,
			Name:         b.Name,
			ServerIDs:    serverIDs,
			Servers:      servers,
			Source:       b.Source,
			Destination:  b.Destination,
			FileType:     b.FileType,
			Type:         b.Type,
			ScheduleType: b.ScheduleType,
			Status:       b.Status,
			SizeBytes:    b.SizeBytes,
			Checksum:     b.Checksum,
			StartedAt:    b.StartedAt,
			CompletedAt:  b.CompletedAt,
			DurationSec:  b.DurationSec,
			ExecutedBy:   b.ExecutedBy,
			CreatedAt:    b.CreatedAt,
			UpdatedAt:    b.UpdatedAt,
		})
	}

	if resp == nil {
		resp = []BackupResponse{}
	}

	return c.JSON(resp)
}
func getBackup(c *fiber.Ctx) error {
	id := c.Params("id")
	var b db.Backup
	if err := db.DB.First(&b, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Backup not found"})
	}

	var serverIDs []uint
	if len(b.ServerIDs) > 0 {
		if err := json.Unmarshal(b.ServerIDs, &serverIDs); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Invalid server_ids format"})
		}
	}

	var servers []db.Server
	if len(serverIDs) > 0 {
		db.DB.Find(&servers, serverIDs)
	}

	resp := BackupResponse{
		ID:           b.ID,
		Name:         b.Name,
		ServerIDs:    serverIDs,
		Servers:      servers,
		Source:       b.Source,
		Destination:  b.Destination,
		FileType:     b.FileType,
		Type:         b.Type,
		ScheduleType: b.ScheduleType,
		Status:       b.Status,
		SizeBytes:    b.SizeBytes,
		Checksum:     b.Checksum,
		StartedAt:    b.StartedAt,
		CompletedAt:  b.CompletedAt,
		DurationSec:  b.DurationSec,
		ExecutedBy:   b.ExecutedBy,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}

	return c.JSON(resp)
}

func createBackup(c *fiber.Ctx) error {
	var backup db.Backup
	if err := c.BodyParser(&backup); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if username := c.Locals("username"); username != nil {
		backup.ExecutedBy = username.(string)
	}

	backup.Status = "pending"
	if err := db.DB.Create(&backup).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(backup)
}

func updateBackup(c *fiber.Ctx) error {
	id := c.Params("id")
	var backup db.Backup
	if err := db.DB.First(&backup, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Backup not found"})
	}

	var updateData db.Backup
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	db.DB.Model(&backup).Updates(updateData)
	return c.JSON(backup)
}

func deleteBackup(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.DB.Delete(&db.Backup{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func executeBackup(c *fiber.Ctx) error {
	id := c.Params("id")
	var backup db.Backup
	if err := db.DB.First(&backup, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Backup not found"})
	}

    // Check if backup is already running
    if backup.Status == "running" {
        // Check if there's an active progress record
        var progress db.BackupProgress
        if err := db.DB.Where("backup_id = ? AND status = ?", backup.ID, "running").First(&progress).Error; err == nil {
            // If the progress hasn't updated recently, treat it as stale and allow restart
            if time.Since(progress.UpdatedAt) > 60*time.Second {
                // Mark existing as failed and reset backup
                progress.Status = "failed"
                progress.Message = "Previous run considered stale; restarting"
                progress.UpdatedAt = time.Now()
                db.DB.Save(&progress)
            } else {
                return c.Status(409).JSON(fiber.Map{"error": "Backup is already running"})
            }
        }
        // No active progress record or stale -> clean old progresses for this backup
        db.DB.Where("backup_id = ?", backup.ID).Delete(&db.BackupProgress{})
        backup.Status = "pending"
        db.DB.Save(&backup)
    }

	// Update status to running and set start time
	now := time.Now()
	backup.Status = "running"
	backup.StartedAt = &now
	if err := db.DB.Save(&backup).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update backup status"})
	}

	// Execute backup using the service
	go func() {
		defer func() {
			// Reset backup status if execution fails catastrophically
			if r := recover(); r != nil {
				backup.Status = "pending"
				db.DB.Save(&backup)
			}
		}()
		backupService.ExecuteBackupAsync(backup)
	}()

	return c.JSON(fiber.Map{"message": "Backup execution started"})
}

func getBackupProgress(c *fiber.Ctx) error {
	id := c.Params("id")
	backupID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid backup ID"})
	}

	progress, err := backupService.GetBackupProgress(uint(backupID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Progress not found"})
	}

	return c.JSON(progress)
}

func getRunningBackups(c *fiber.Ctx) error {
	progresses, err := backupService.GetAllRunningBackups()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get running backups"})
	}

	// Convert to response format with backup data
	responses := []BackupProgressResponse{}
	for _, progress := range progresses {
		// Load backup data
		var backup db.Backup
		if err := db.DB.First(&backup, progress.BackupID).Error; err != nil {
			continue // Skip if backup not found
		}

		// Load server data
		var serverIDs []uint
		if len(backup.ServerIDs) > 0 {
			if err := json.Unmarshal(backup.ServerIDs, &serverIDs); err != nil {
				serverIDs = []uint{}
			}
		}
		var servers []db.Server
		if len(serverIDs) > 0 {
			db.DB.Find(&servers, serverIDs)
		}

		responses = append(responses, BackupProgressResponse{
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
		})
	}

	return c.JSON(responses)
}

func deleteAllProcesses(c *fiber.Ctx) error {
	// Delete all backup progress records that are running, failed, or completed
	if err := db.DB.Where("status IN ?", []string{"running", "failed", "completed"}).Delete(&db.BackupProgress{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func deleteProcess(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.DB.Delete(&db.BackupProgress{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

