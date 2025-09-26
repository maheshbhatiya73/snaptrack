package db

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Server struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;unique" json:"name"`
	Host        string         `gorm:"not null" json:"host"`
	SSHUser     *string        `json:"ssh_user"`
	SSHPort     *int           `json:"ssh_port"`
	SSHKeyPath  *string        `json:"ssh_key_path"`
	Type        string         `gorm:"not null" json:"type"` // local / remote
	TransferType *string `json:"transferType"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Backup struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null;uniqueIndex" json:"name"`
	Source       string         `gorm:"not null" json:"source"`
	Destination  string         `gorm:"not null" json:"destination"`
	FileType     string         `gorm:"not null" json:"file_type"` // tar / zip / raw
	Type         string         `gorm:"not null" json:"type"`      // full / incremental
	ScheduleType string         `gorm:"not null;default:one_time" json:"schedule_type"`
	Status       string         `gorm:"not null" json:"status"`    // scheduled / running / success / failed
	ServerIDs    datatypes.JSON `gorm:"type:jsonb;not null" json:"server_ids"`
	
	SizeBytes    int64          `json:"size_bytes"`
	Checksum     *string        `json:"checksum"`
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	DurationSec  int64          `json:"duration_sec"`
	ExecutedBy   string         `gorm:"not null" json:"executed_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Log struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Level     string         `gorm:"not null" json:"level"`   // info / warning / error
	Message   string         `gorm:"not null" json:"message"`
	Context   *string        `json:"context"`                 // optional string
	EntityType *string       `json:"entity_type"`             // "backup", "server", etc.
	EntityID   *uint         `json:"entity_id"`               // optional reference to a backup/server
	Metadata   datatypes.JSON `json:"metadata"`               // store extra data (JSON)
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


type BackupProgress struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	BackupID    uint      `gorm:"not null;index" json:"backup_id"`
	Backup      Backup    `gorm:"foreignKey:BackupID" json:"-"`
	Status      string    `gorm:"not null" json:"status"` // pending / running / completed / failed
	Progress    int       `gorm:"default:0" json:"progress"` // 0-100
	Message     string    `gorm:"not null" json:"message"`
	CurrentFile *string   `json:"current_file"`
	BytesProcessed int64  `gorm:"default:0" json:"bytes_processed"`
	TotalBytes    *int64  `json:"total_bytes"`
	SpeedBPS     *int64   `json:"speed_bps"` // bytes per second
	ETASeconds  *int64   `json:"eta_seconds"` // estimated time remaining
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
