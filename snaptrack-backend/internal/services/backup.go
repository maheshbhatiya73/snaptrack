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
	ctx := s.cron.Stop()
	<-ctx.Done() // Wait for all running jobs to complete
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

	// Find backups that are pending or completed (for recurring backups)
	cursor, err := collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"status": models.StatusPending},
			{"status": models.StatusCompleted, "schedule.kind": bson.M{"$ne": models.ScheduleOneTime}},
		},
	})
	if err != nil {
		log.Printf("Error fetching backups: %v", err)
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

func (s *BackupService) scheduleBackup(backup models.Backup) {
	// Skip if backup is already completed or failed for one-time schedule
	if backup.Schedule.Kind == models.ScheduleOneTime && (backup.Status == models.StatusCompleted || backup.Status == models.StatusFailed) {
		log.Printf("Skipping one-time backup %s: already %s", backup.ID.Hex(), backup.Status)
		return
	}

	var spec string
	now := time.Now().UTC()
	scheduleTime := backup.Schedule.Date

	// Update NextRun for recurring backups or one-time future backups
	var nextRun time.Time
	switch backup.Schedule.Kind {
	case models.ScheduleOneTime:
		if scheduleTime.After(now) {
			spec = scheduleTime.Format("15 04 02 01 * 2006")
			nextRun = scheduleTime
		} else {
			s.executeBackup(backup)
			return
		}
	case models.ScheduleHourly:
		spec = "0 * * * *"
		nextRun = now.Truncate(time.Hour).Add(time.Hour)
	case models.ScheduleDaily:
		spec = "0 0 * * *"
		nextRun = now.Truncate(24 * time.Hour).Add(24 * time.Hour)
	case models.ScheduleWeekly:
		spec = "0 0 * * 0"
		nextRun = now.Truncate(24 * time.Hour).AddDate(0, 0, 7-int(now.Weekday()))
	case models.ScheduleMonthly:
		spec = "0 0 1 * *"
		nextRun = now.Truncate(24 * time.Hour).AddDate(0, 1, 1-int(now.Day()))
	default:
		log.Printf("Invalid schedule kind for backup %s: %s", backup.ID.Hex(), backup.Schedule.Kind)
		s.createBackupLog(backup.ID, "failed", "Invalid schedule kind")
		return
	}

	// Check if backup is already scheduled by looking for a "scheduled" log
	collection := s.db.Collection("backups")
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	var currentBackup models.Backup
	err := collection.FindOne(ctx, bson.M{
		"_id": backup.ID,
		"logs": bson.M{
			"$elemMatch": bson.M{
				"status": "scheduled",
				"message": bson.M{"$regex": fmt.Sprintf("Scheduled backup %s", backup.ID.Hex())},
			},
		},
	}).Decode(&currentBackup)
	isScheduled := err == nil

	// Update NextRun in the database if not already scheduled or if outdated
	if !isScheduled || backup.NextRun.IsZero() || backup.NextRun.Before(now) {
		_, err = collection.UpdateOne(ctx, bson.M{"_id": backup.ID}, bson.M{
			"$set": bson.M{"nextRun": nextRun},
		})
		if err != nil {
			log.Printf("Error updating nextRun for backup %s: %v", backup.ID.Hex(), err)
			s.createBackupLog(backup.ID, "failed", "Failed to update nextRun")
			return
		}
	}

	// Log scheduling only if not already scheduled
	if !isScheduled {
		countdown := time.Until(nextRun).Round(time.Second).String()
		logMsg := fmt.Sprintf("Scheduled backup %s (%s) for %s (in %s)", backup.ID.Hex(), backup.Schedule.Kind, nextRun.Format(time.RFC3339), countdown)
		log.Println(logMsg)
		s.createBackupLog(backup.ID, "scheduled", logMsg)
	}

	// Add to cron only if not already scheduled
	if !isScheduled {
		_, err := s.cron.AddFunc(spec, func() {
			s.executeBackup(backup)
			if backup.Schedule.Kind != models.ScheduleOneTime {
				s.updateNextRun(backup)
			}
		})
		if err != nil {
			log.Printf("Error scheduling backup %s: %v", backup.ID.Hex(), err)
			s.createBackupLog(backup.ID, "failed", "Failed to schedule backup")
			return
		}
	}
}

// updateNextRun updates the nextRun timestamp for recurring backups
func (s *BackupService) updateNextRun(backup models.Backup) {
	now := time.Now().UTC()
	var nextRun time.Time

	switch backup.Schedule.Kind {
	case models.ScheduleHourly:
		nextRun = now.Truncate(time.Hour).Add(time.Hour)
	case models.ScheduleDaily:
		nextRun = now.Truncate(24 * time.Hour).Add(24 * time.Hour)
	case models.ScheduleWeekly:
		nextRun = now.Truncate(24 * time.Hour).AddDate(0, 0, 7-int(now.Weekday()))
	case models.ScheduleMonthly:
		nextRun = now.Truncate(24 * time.Hour).AddDate(0, 1, 1-int(now.Day()))
	default:
		return // Should not happen
	}

	_, err := s.db.Collection("backups").UpdateOne(s.ctx, bson.M{"_id": backup.ID}, bson.M{
		"$set": bson.M{"nextRun": nextRun},
	})
	if err != nil {
		log.Printf("Error updating nextRun for backup %s: %v", backup.ID.Hex(), err)
		s.createBackupLog(backup.ID, "failed", "Failed to update nextRun")
		return
	}

	// Do not create a new "scheduled" log here to avoid duplicates
	log.Printf("Updated nextRun for backup %s (%s) to %s", backup.ID.Hex(), backup.Schedule.Kind, nextRun.Format(time.RFC3339))
}

// executeBackup performs the backup operation
func (s *BackupService) executeBackup(backup models.Backup) {
	collection := s.db.Collection("backups")
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Minute)
	defer cancel()

	// Check if backup is already running to prevent duplicates
	var currentBackup models.Backup
	err := collection.FindOne(ctx, bson.M{"_id": backup.ID}).Decode(&currentBackup)
	if err != nil {
		log.Printf("Error checking backup %s status: %v", backup.ID.Hex(), err)
		s.createBackupLog(backup.ID, "failed", "Failed to check backup status")
		return
	}
	if currentBackup.Status == models.StatusStarted {
		log.Printf("Backup %s is already running, skipping execution", backup.ID.Hex())
		s.createBackupLog(backup.ID, "skipped", "Backup is already running")
		return
	}

	// Update status to started
	_, err = collection.UpdateOne(ctx, bson.M{"_id": backup.ID}, bson.M{
		"$set": bson.M{"status": models.StatusStarted},
	})
	if err != nil {
		log.Printf("Error updating backup status to started: %v", err)
		s.createBackupLog(backup.ID, "failed", "Failed to update status to started")
		return
	}

	s.createBackupLog(backup.ID, "started", fmt.Sprintf("Backup started: %s to %s", backup.SourcePath, backup.DestinationPath))
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
	if err != nil {
		log.Printf("Backup %s failed: %v", backup.ID.Hex(), err)
		s.createBackupLog(backup.ID, "failed", fmt.Sprintf("Backup failed: %s", string(output)))
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
