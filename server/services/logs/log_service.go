package logs

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Log struct {
	ID        uint           `gorm:"primaryKey"`
	Level     string
	Message   string
	EntityType *string
	EntityID   *uint
	Metadata   json.RawMessage
	CreatedAt  time.Time
}

type LogService struct {
	DB *gorm.DB
}

func NewLogService(db *gorm.DB) *LogService {
	return &LogService{DB: db}
}

func (ls *LogService) Insert(level, message string, entityType *string, entityID *uint, metadata map[string]interface{}) {
	metadataJSON, _ := json.Marshal(metadata)

	log := Log{
		Level:      level,
		Message:    message,
		EntityType: entityType,
		EntityID:   entityID,
		Metadata:   metadataJSON,
		CreatedAt:  time.Now(),
	}

	ls.DB.Create(&log)
}

func (ls *LogService) Info(message string, entityType *string, entityID *uint, metadata map[string]interface{}) {
	ls.Insert("info", message, entityType, entityID, metadata)
}

func PtrString(s string) *string { return &s }
func PtrUint(u uint) *uint       { return &u }
