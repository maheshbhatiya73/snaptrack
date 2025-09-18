package db

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;unique" json:"name"`
	Host        string         `gorm:"not null" json:"host"`
	SSHUser     *string        `json:"ssh_user"`     // nullable for local servers
	SSHPort     *int           `json:"ssh_port"`     // nullable for local servers
	SSHKeyPath  *string        `json:"ssh_key_path"` // nullable for local servers
	Type        string         `gorm:"not null" json:"type"` // "local" or "remote"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Backups     []Backup       `json:"backups,omitempty"`
}

type Backup struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ServerID    uint           `gorm:"index" json:"server_id"`
	Server      Server         `json:"server"`
	Name        string         `gorm:"not null;uniqueIndex" json:"name"`
	Source      string         `gorm:"not null" json:"source"`
	Destination string         `gorm:"not null" json:"destination"`
	SizeBytes   int64          `json:"size_bytes"`
	Checksum    *string        `json:"checksum"`
	Status      string         `gorm:"not null" json:"status"` // success / failed
	Type        string         `gorm:"not null" json:"type"`   // full / incremental
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt time.Time      `json:"completed_at"`
	DurationSec int64          `json:"duration_sec"`
	ExecutedBy  string         `gorm:"not null" json:"executed_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Example Log table (optional but recommended)
type Log struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Level     string         `gorm:"not null" json:"level"` // info, warn, error
	Message   string         `gorm:"not null" json:"message"`
	Context   *string        `json:"context"`               // optional (e.g. backup id)
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
