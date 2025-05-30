package services

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var BackupCollection *mongo.Collection

func InitBackupService(db *mongo.Database) {
	BackupCollection = db.Collection("backups")
	go StartBackupProcessor() // 🔁 Start background job loop
	fmt.Println("✅ Backup processor started")
}

type Schedule struct {
	Kind string    `bson:"kind" json:"kind"`
	Date time.Time `bson:"date" json:"date"`
}

type BackupJob struct {
	ID              string    `bson:"_id" json:"id"`
	App             string    `bson:"app" json:"app"`
	Type            string    `bson:"type" json:"type"`
	Size            string    `bson:"size" json:"size"`
	Status          string    `bson:"status" json:"status"`
	SourcePath      string    `bson:"sourcePath" json:"sourcePath"`
	DestinationPath string    `bson:"destinationPath" json:"destinationPath"`
	FileType        string    `bson:"fileType" json:"fileType"`
	Schedule        Schedule  `bson:"schedule" json:"schedule"`
	CreatedAt       time.Time `bson:"createdAt" json:"createdAt"`
}

func StartBackupProcessor() {
	for {
		time.Sleep(10 * time.Second)

		backups, err := GetPendingBackups()
		if err != nil {
			fmt.Println("❌ Failed to fetch pending backups:", err)
			continue
		}

		for _, job := range backups {
			if job.Schedule.Kind == "one-time" && time.Now().Before(job.Schedule.Date) {
				continue
			}
			fmt.Println("📦 Running backup for:", job.App)
			err := RunBackup(job)
			if err != nil {
				fmt.Println("❌ Backup failed:", err)
				UpdateBackupStatus(job.ID, "failed")
			} else {
				UpdateBackupStatus(job.ID, "completed")
			}
		}
	}
}

func GetPendingBackups() ([]BackupJob, error) {
	var results []BackupJob
	filter := bson.M{"status": "pending"}
	cursor, err := BackupCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func UpdateBackupStatus(id string, status string) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := BackupCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("⚠️ Failed to update status:", err)
	}
}
func RunBackup(job BackupJob) error {
	source := job.SourcePath
	destDir := job.DestinationPath
	filename := fmt.Sprintf("%s-%d.tar.gz", job.App, time.Now().Unix())

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	outFilePath := filepath.Join(destDir, filename)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer outFile.Close()

	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	err = filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, file)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tarWriter, f)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed during file walk: %v", err)
	}

	fmt.Println("✅ Backup completed:", outFilePath)
	return nil
}

func PrintAllBackupJobs() {
	backups, err := GetPendingBackups()
	if err != nil {
		fmt.Println("❌ Failed to fetch backups:", err)
		return
	}

	jsonData, err := json.MarshalIndent(backups, "", "  ")
	if err != nil {
		fmt.Println("❌ Failed to marshal backups:", err)
		return
	}

	fmt.Println("📦 All Pending Backups:")
	fmt.Println(string(jsonData))
}
