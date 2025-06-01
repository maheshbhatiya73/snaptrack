package controller

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"fmt"
	"snaptrackserver/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrorResponse defines the structure for JSON error responses
type ErrorResponse struct {
	Message string `json:"message"`
}

var collection *mongo.Collection

// SetCollection initializes the MongoDB collection
func SetCollection(c *mongo.Collection) {
	collection = c
	log.Println("Backup collection initialized")
}

// writeJSONError writes a JSON error response with the given status code and message
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Message: message}); err != nil {
		log.Printf("writeJSONError: Failed to encode error response: %v", err)
	}
}

// getFolderSize calculates the total size of all files in the directory (recursive)
func getFolderSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return size, nil
}

// formatSize converts bytes to a human-readable format (e.g., "1.2MB")
func formatSize(bytes int64) string {
	if bytes == 0 {
		return "0B"
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(bytes)
	unitIndex := 0
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}
	return fmt.Sprintf("%.1f%s", size, units[unitIndex])
}

func CreateBackup(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("CreateBackup: Database not initialized")
		writeJSONError(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("CreateBackup: Invalid method %s", r.Method)
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var backup models.Backup
	if err := json.NewDecoder(r.Body).Decode(&backup); err != nil {
		log.Printf("CreateBackup: Failed to decode request body: %v", err)
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate SourcePath
	if _, err := os.Stat(backup.SourcePath); os.IsNotExist(err) {
		log.Printf("CreateBackup: Source path does not exist: %s", backup.SourcePath)
		writeJSONError(w, "Source path does not exist", http.StatusBadRequest)
		return
	}

	// Validate DestinationPath
	if _, err := os.Stat(backup.DestinationPath); os.IsNotExist(err) {
		log.Printf("CreateBackup: Destination path does not exist: %s", backup.DestinationPath)
		writeJSONError(w, "Destination path does not exist", http.StatusBadRequest)
		return
	}

	// Calculate source folder size
	sizeBytes, err := getFolderSize(backup.SourcePath)
	if err != nil {
		log.Printf("CreateBackup: Failed to calculate source folder size for %s: %v", backup.SourcePath, err)
		writeJSONError(w, "Failed to calculate source folder size", http.StatusInternalServerError)
		return
	}
	backup.Size = formatSize(sizeBytes)

	backup.Status = models.StatusPending // Set initial status
	backup.CreatedAt = time.Now()
	backup.ID = primitive.NewObjectID()

	res, err := collection.InsertOne(context.Background(), backup)
	if err != nil {
		log.Printf("CreateBackup: Failed to insert backup: %v", err)
		writeJSONError(w, "Failed to create backup", http.StatusInternalServerError)
		return
	}

	log.Printf("CreateBackup: Backup created with ID %v, Size %s", res.InsertedID, backup.Size)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(backup); err != nil {
		log.Printf("CreateBackup: Failed to encode response: %v", err)
	}
}

func GetAllBackups(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("GetAllBackups: Database not initialized")
		writeJSONError(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("GetAllBackups: Failed to fetch backups: %v", err)
		writeJSONError(w, "Failed to fetch backups", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var backups []models.Backup
	if err := cursor.All(context.Background(), &backups); err != nil {
		log.Printf("GetAllBackups: Failed to decode backups: %v", err)
		writeJSONError(w, "Failed to decode backups", http.StatusInternalServerError)
		return
	}

	log.Printf("GetAllBackups: Returned %d backups", len(backups))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(backups); err != nil {
		log.Printf("GetAllBackups: Failed to encode response: %v", err)
	}
}

func GetBackupByID(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("GetBackupByID: Database not initialized")
		writeJSONError(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/api/backups/"):]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("GetBackupByID: Invalid backup ID %s", idStr)
		writeJSONError(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	var backup models.Backup
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&backup)
	if err != nil {
		log.Printf("GetBackupByID: Backup not found with ID %s", idStr)
		writeJSONError(w, "Backup not found", http.StatusNotFound)
		return
	}

	log.Printf("GetBackupByID: Returned backup with ID %s", idStr)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(backup); err != nil {
		log.Printf("GetBackupByID: Failed to encode response: %v", err)
	}
}

func UpdateBackup(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("UpdateBackup: Database not initialized")
		writeJSONError(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/api/backups/"):]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("UpdateBackup: Invalid backup ID %s", idStr)
		writeJSONError(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	var backup models.Backup
	if err := json.NewDecoder(r.Body).Decode(&backup); err != nil {
		log.Printf("UpdateBackup: Failed to decode request body: %v", err)
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": backup}
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		log.Printf("UpdateBackup: Failed to update backup with ID %s: %v", idStr, err)
		writeJSONError(w, "Failed to update backup", http.StatusInternalServerError)
		return
	}

	if res.MatchedCount == 0 {
		log.Printf("UpdateBackup: No backup found to update with ID %s", idStr)
		writeJSONError(w, "Backup not found", http.StatusNotFound)
		return
	}

	log.Printf("UpdateBackup: Backup updated with ID %s", idStr)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Backup updated successfully"}); err != nil {
		log.Printf("UpdateBackup: Failed to encode response: %v", err)
	}
}

func DeleteBackup(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("DeleteBackup: Database not initialized")
		writeJSONError(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/api/backups/"):]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("DeleteBackup: Invalid backup ID %s", idStr)
		writeJSONError(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Printf("DeleteBackup: Failed to delete backup with ID %s: %v", idStr, err)
		writeJSONError(w, "Failed to delete backup", http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		log.Printf("DeleteBackup: No backup found to delete with ID %s", idStr)
		writeJSONError(w, "Backup not found", http.StatusNotFound)
		return
	}

	log.Printf("DeleteBackup: Backup deleted with ID %s", idStr)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Backup deleted successfully"}); err != nil {
		log.Printf("DeleteBackup: Failed to encode response: %v", err)
	}
}
