package db

import "time"

type Backup struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Source      string
	Destination string
	Size        int64
	Status      string
	Type        string
	Checksum    string
	CreatedAt   time.Time
}

type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Role      string
	LastLogin time.Time
}

type Log struct {
	ID        uint `gorm:"primaryKey"`
	BackupID  uint
	Message   string
	Level     string
	CreatedAt time.Time
}
