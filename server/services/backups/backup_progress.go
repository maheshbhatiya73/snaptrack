package backups

import (
    "snaptrack/db"
    "time"
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
