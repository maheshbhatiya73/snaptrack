package routes

import (
	"encoding/json"
	"snaptrack/db"

	"github.com/gofiber/fiber/v2"
)

// LogResponse represents the structured JSON log format
type LogResponse struct {
	ID         uint        `json:"id"`
	Level      string      `json:"level"`
	Message    string      `json:"message"`
	EntityType *string     `json:"entity_type"`
	EntityID   *uint       `json:"entity_id"`
	EntityName *string     `json:"entity_name"` // new field
	Metadata   interface{} `json:"metadata"`
	CreatedAt  string      `json:"created_at"`
}

func RecentActivity(c *fiber.Ctx) error {
	var logs []db.Log
	var response []LogResponse

	if err := db.DB.Order("created_at desc").Limit(10).Find(&logs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch recent activity"})
	}

	for _, log := range logs {
		var parsedMetadata interface{}
		if log.Metadata != nil {
			if err := json.Unmarshal(log.Metadata, &parsedMetadata); err != nil {
				parsedMetadata = log.Metadata
			}
		}

		var entityName *string
		if log.EntityType != nil && log.EntityID != nil {
			switch *log.EntityType {
			case "backup":
				var backup db.Backup
				if err := db.DB.First(&backup, *log.EntityID).Error; err == nil {
					entityName = &backup.Name
				}
			case "server":
				var server db.Server
				if err := db.DB.First(&server, *log.EntityID).Error; err == nil {
					entityName = &server.Name
				}
			}
		}

		response = append(response, LogResponse{
			ID:         log.ID,
			Level:      log.Level,
			Message:    log.Message,
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			EntityName: entityName,
			Metadata:   parsedMetadata,
			CreatedAt:  log.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return c.JSON(response)
}
