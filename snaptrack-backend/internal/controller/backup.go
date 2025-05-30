package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"snaptrackserver/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func SetCollection(c *mongo.Collection) {
	collection = c
	log.Println("Backup collection initialized")
}

func CreateBackup(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("CreateBackup: Database not initialized")
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("CreateBackup: Invalid method %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var backup models.Backup
	if err := json.NewDecoder(r.Body).Decode(&backup); err != nil {
		log.Printf("CreateBackup: Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	backup.CreatedAt = time.Now()
	backup.ID = primitive.NewObjectID()

	res, err := collection.InsertOne(context.Background(), backup)
	if err != nil {
		log.Printf("CreateBackup: Failed to insert backup: %v\n", err)
		http.Error(w, "Failed to create backup", http.StatusInternalServerError)
		return
	}

	log.Printf("CreateBackup: Backup created with ID %v\n", res.InsertedID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(backup)
}


func GetAllBackups(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("GetAllBackups: Database not initialized")
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("GetAllBackups: Failed to fetch backups: %v\n", err)
		http.Error(w, "Failed to fetch backups", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var backups []models.Backup
	if err := cursor.All(context.Background(), &backups); err != nil {
		log.Printf("GetAllBackups: Failed to decode backups: %v\n", err)
		http.Error(w, "Failed to decode backups", http.StatusInternalServerError)
		return
	}

	log.Printf("GetAllBackups: Returned %d backups\n", len(backups))
	json.NewEncoder(w).Encode(backups)
}

func GetBackupByID(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("GetBackupByID: Database not initialized")
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/api/backups/"):]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("GetBackupByID: Invalid backup ID %s\n", idStr)
		http.Error(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	var backup models.Backup
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&backup)
	if err != nil {
		log.Printf("GetBackupByID: Backup not found with ID %s\n", idStr)
		http.Error(w, "Backup not found", http.StatusNotFound)
		return
	}

	log.Printf("GetBackupByID: Returned backup with ID %s\n", idStr)
	json.NewEncoder(w).Encode(backup)
}

func UpdateBackup(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("UpdateBackup: Database not initialized")
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/api/backups/"):]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("UpdateBackup: Invalid backup ID %s\n", idStr)
		http.Error(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	var backup models.Backup
	if err := json.NewDecoder(r.Body).Decode(&backup); err != nil {
		log.Printf("UpdateBackup: Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": backup}
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		log.Printf("UpdateBackup: Failed to update backup with ID %s: %v\n", idStr, err)
		http.Error(w, "Failed to update backup", http.StatusInternalServerError)
		return
	}

	if res.MatchedCount == 0 {
		log.Printf("UpdateBackup: No backup found to update with ID %s\n", idStr)
		http.Error(w, "Backup not found", http.StatusNotFound)
		return
	}

	log.Printf("UpdateBackup: Backup updated with ID %s\n", idStr)
	json.NewEncoder(w).Encode(map[string]string{"message": "Backup updated successfully"})
}


func DeleteBackup(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		log.Println("DeleteBackup: Database not initialized")
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Path[len("/api/backups/"):]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("DeleteBackup: Invalid backup ID %s\n", idStr)
		http.Error(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Printf("DeleteBackup: Failed to delete backup with ID %s: %v\n", idStr, err)
		http.Error(w, "Failed to delete backup", http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		log.Printf("DeleteBackup: No backup found to delete with ID %s\n", idStr)
		http.Error(w, "Backup not found", http.StatusNotFound)
		return
	}

	log.Printf("DeleteBackup: Backup deleted with ID %s\n", idStr)
	json.NewEncoder(w).Encode(map[string]string{"message": "Backup deleted successfully"})
}
