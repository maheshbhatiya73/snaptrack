package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"snaptrackserver/internal/models"
)

// BackupService manages backup operations and scheduling
type BackupService struct {
	db   *mongo.Database
	cron *cron.Cron
	ctx  context.Context
}

// backupService is a singleton instance of BackupService
var backupService *BackupService

// InitBackupService initializes and starts the backup service
func InitBackupService(db *mongo.Database) {
	if backupService == nil {
		backupService = NewBackupService(db)
		backupService.Start()
	}
}

// NewBackupService creates a new BackupService instance
func NewBackupService(db *mongo.Database) *BackupService {
	ctx := context.Background()
	svc := &BackupService{
		db:   db,
		cron: cron.New(),
		ctx:  ctx,
	}
	return svc
}

// Start begins the backup scheduler
func (s *BackupService) Start() {
	s.cron.AddFunc("@every 1m", s.checkAndScheduleBackups)
	s.cron.Start()
	log.Println("Backup scheduler started")
}

// Stop gracefully stops the backup scheduler
func (s *BackupService) Stop() {
	s.cron.Stop()
	log.Println("Backup scheduler stopped")
}

// PrintAllBackupJobs prints all scheduled backup jobs
func (s *BackupService) PrintAllBackupJobs() {
	entries := s.cron.Entries()
	if len(entries) == 0 {
		fmt.Println("No backup jobs scheduled.")
		return
	}

	fmt.Println("Scheduled Backup Jobs:")
	for i, entry := range entries {
		fmt.Printf("Job %d:\n", i+1)
		fmt.Printf("  ID: %d\n", entry.ID)
		fmt.Printf("  Schedule: %s\n", entry.Schedule)
		fmt.Printf("  Next Run: %s\n", entry.Next.Format(time.RFC3339))
		fmt.Printf("  Previous Run: %s\n", entry.Prev.Format(time.RFC3339))
	}
}

// checkAndScheduleBackups checks the database for backups and schedules them
func (s *BackupService) checkAndScheduleBackups() {
	collection := s.db.Collection("backups")
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"status": models.StatusPending})
	if err != nil {
		log.Printf("Error fetching pending backups: %v", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var backup models.Backup
		if err := cursor.Decode(&backup); err != nil {
			log.Printf("Error decoding backup: %v", err)
			continue
		}
		s.scheduleBackup(backup)
	}
}

// scheduleBackup schedules a backup based on its schedule kind
func (s *BackupService) scheduleBackup(backup models.Backup) {
	var spec string
	scheduleTime := backup.Schedule.Date
	now := time.Now().UTC()

	switch backup.Schedule.Kind {
	case models.ScheduleOneTime:
		if scheduleTime.After(now) {
			spec = scheduleTime.Format("15 04 02 01 * 2006") // cron format: mm hh DD MM * YYYY
		} else {
			s.executeBackup(backup)
			return
		}
	case models.ScheduleHourly:
		spec = "0 * * * *" // Every hour at minute 0
	case models.ScheduleDaily:
		spec = "0 0 * * *" // Every day at midnight
	case models.ScheduleWeekly:
		spec = "0 0 * * 0" // Every Sunday at midnight
	case models.ScheduleMonthly:
		spec = "0 0 1 * *" // First day of every month at midnight
	default:
		log.Printf("Invalid schedule kind for backup %s: %s", backup.ID.Hex(), backup.Schedule.Kind)
		s.createBackupLog(backup.ID, "failed", "Invalid schedule kind")
		return
	}

	entryID, err := s.cron.AddFunc(spec, func() {
		s.executeBackup(backup)
	})
	if err != nil {
		log.Printf("Error scheduling backup %s: %v", backup.ID.Hex(), err)
		s.createBackupLog(backup.ID, "failed", "Failed to schedule backup")
		return
	}
	log.Printf("Scheduled backup %s with cron entry ID %d", backup.ID.Hex(), entryID)
}

// executeBackup performs the backup operation
func (s *BackupService) executeBackup(backup models.Backup) {
	collection := s.db.Collection("backups")
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Minute)
	defer cancel()

	// Update status to started
	_, err := collection.UpdateOne(ctx, bson.M{"_id": backup.ID}, bson.M{
		"$set": bson.M{"status": models.StatusStarted},
	})
	if err != nil {
		log.Printf("Error updating backup status to started: %v", err)
		s.createBackupLog(backup.ID, "failed", "Failed to update status to started")
		return
	}

	s.createBackupLog(backup.ID, "started", "Backup operation started")
	log.Printf("Backup started: ID=%s, App=%s", backup.ID.Hex(), backup.App)

	// Construct backup command based on FileType
	var cmd *exec.Cmd
	filename := backup.App + "-" + time.Now().UTC().Format("20060102T150405") + "." + string(backup.FileType)
	destination := filepath.Join(backup.DestinationPath, filename)

	switch backup.FileType {
	case models.FileTypeZIP:
		cmd = exec.Command("zip", "-r", destination, backup.SourcePath)
	case models.FileTypeTAR:
		cmd = exec.Command("tar", "-cvf", destination, backup.SourcePath)
	case models.FileTypeTARGZ:
		cmd = exec.Command("tar", "-zcvf", destination, backup.SourcePath)
	default:
		log.Printf("Unsupported file type for backup %s: %s", backup.ID.Hex(), backup.FileType)
		s.createBackupLog(backup.ID, "failed", "Unsupported file type")
		return
	}

	// Execute backup command
	output, err := cmd.CombinedOutput()
	s.createBackupLog(backup.ID, "executed", "Backup command executed")
	if err != nil {
		log.Printf("Backup %s failed: %v", backup.ID.Hex(), err)
		s.createBackupLog(backup.ID, "failed", "Backup failed: "+string(output))
		// Update backup status to failed
		_, err = collection.UpdateOne(ctx, bson.M{"_id": backup.ID}, bson.M{
			"$set": bson.M{
				"status": models.StatusFailed,
				"size":   s.getFileSize(destination),
			},
		})
		if err != nil {
			log.Printf("Error updating backup status: %v", err)
			s.createBackupLog(backup.ID, "failed", "Failed to update backup status")
		}
		return
	}

	// Update backup status to completed
	_, err = collection.UpdateOne(ctx, bson.M{"_id": backup.ID}, bson.M{
		"$set": bson.M{
			"status": models.StatusCompleted,
			"size":   s.getFileSize(destination),
		},
	})
	if err != nil {
		log.Printf("Error updating backup status: %v", err)
		s.createBackupLog(backup.ID, "failed", "Failed to update backup status")
		return
	}

	s.createBackupLog(backup.ID, "completed", "Backup completed successfully")
}

// createBackupLog appends a log entry to the backup's Logs field
func (s *BackupService) createBackupLog(backupID primitive.ObjectID, status, message string) {
	collection := s.db.Collection("backups")
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	logEntry := models.BackupLog{
		ID:        primitive.NewObjectID(),
		BackupID:  backupID,
		Status:    status,
		Message:   message,
		CreatedAt: time.Now().UTC(),
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": backupID}, bson.M{
		"$push": bson.M{"logs": logEntry},
	})
	if err != nil {
		log.Printf("Error appending backup log: %v", err)
	}
}

// getFileSize returns the size of a file in human-readable format
func (s *BackupService) getFileSize(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "0B"
	}
	bytes := fileInfo.Size()
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(bytes)
	unitIndex := 0
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}
	return fmt.Sprintf("%.1f%s", size, units[unitIndex])
}