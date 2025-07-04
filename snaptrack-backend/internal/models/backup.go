package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BackupType defines valid types of backups
type BackupType string

const (
	BackupTypeManual      BackupType = "manual"
	BackupTypeFull        BackupType = "full"
	BackupTypeIncremental BackupType = "incremental"
)

// FileType defines valid types of backup file formats
type FileType string

const (
	FileTypeZIP    FileType = "zip"
	FileTypeTAR    FileType = "tar"
	FileTypeTARGZ  FileType = "tar.gz"
)

// ScheduleKind defines the type of scheduling
type ScheduleKind string

const (
	ScheduleOneTime ScheduleKind = "one-time"
	ScheduleHourly  ScheduleKind = "hourly"
	ScheduleDaily   ScheduleKind = "daily"
	ScheduleWeekly  ScheduleKind = "weekly"
	ScheduleMonthly ScheduleKind = "monthly"
)

// Schedule structure
type Schedule struct {
	Kind ScheduleKind `json:"kind" bson:"kind"`
	Date time.Time    `json:"date" bson:"date"`
}

type BackupStatus string

const (
	StatusPending   BackupStatus = "pending"
	StatusStarted   BackupStatus = "started"
	StatusCompleted BackupStatus = "completed"
	StatusFailed    BackupStatus = "failed"
)

// BackupLog represents a log entry for a backup operation
type BackupLog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BackupID  primitive.ObjectID `json:"backupId" bson:"backupId"` // Reference to Backup
	Status    string             `json:"status" bson:"status"`     // e.g., "started", "executed", "completed", "failed"
	Message   string             `json:"message" bson:"message"`   // Log message or error details
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`

}

// Backup model
type Backup struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	App             string             `json:"app" bson:"app"`
	Type            BackupType         `json:"type" bson:"type"`
	Size            string             `json:"size,omitempty" bson:"size,omitempty"`
	Status          BackupStatus       `json:"status,omitempty" bson:"status,omitempty"`
	SourcePath      string             `json:"sourcePath" bson:"sourcePath"`
	DestinationPath string             `json:"destinationPath" bson:"destinationPath"`
	FileType        FileType           `json:"fileType" bson:"fileType"`
	Schedule        Schedule           `json:"schedule" bson:"schedule"`
	NextRun         time.Time          `json:"nextRun,omitempty" bson:"nextRun,omitempty"`
	Logs            []BackupLog        `json:"logs,omitempty" bson:"logs,omitempty"`
	RunNow          bool               `json:"runNow" bson:"runNow"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
}