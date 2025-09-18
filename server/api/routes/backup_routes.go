package routes

import (
	"encoding/json"
	"snaptrack/auth"
	"snaptrack/db"
	"time"

	"github.com/gofiber/fiber/v2"
)

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

func RegisterBackupRoutes(app *fiber.App) {
	api := app.Group("/api/backups", auth.RequireJWT())

	api.Get("/", listBackups)
	api.Get("/:id", getBackup)
	api.Post("/", createBackup)
	api.Put("/:id", updateBackup)
	api.Delete("/:id", deleteBackup)
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
			StartedAt:    &b.StartedAt,
			CompletedAt:  &b.CompletedAt,
			DurationSec:  b.DurationSec,
			ExecutedBy:   b.ExecutedBy,
			CreatedAt:    b.CreatedAt,
			UpdatedAt:    b.UpdatedAt,
		})
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
		StartedAt:    &b.StartedAt,
		CompletedAt:  &b.CompletedAt,
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
