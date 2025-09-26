package db

import (
	"snaptrack/services/logs"
	"gorm.io/gorm"
)

func (b *Backup) AfterCreate(tx *gorm.DB) (err error) {
	ls := logs.NewLogService(tx)
	ls.Info("Backup created", logs.PtrString("backup"), &b.ID, map[string]interface{}{
		"backup_name": b.Name,
	})
	return
}

func (b *Backup) AfterUpdate(tx *gorm.DB) (err error) {
	ls := logs.NewLogService(tx)
	ls.Info("Backup updated", logs.PtrString("backup"), &b.ID, map[string]interface{}{
		"status": b.Status,
	})
	return
}

func (b *Backup) AfterDelete(tx *gorm.DB) (err error) {
	ls := logs.NewLogService(tx)
	ls.Info("Backup deleted", logs.PtrString("backup"), &b.ID, nil)
	return
}
