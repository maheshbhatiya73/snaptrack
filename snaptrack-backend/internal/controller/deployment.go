package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"snaptrackserver/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var deploymentCollection *mongo.Collection

func SetDeploymentCollection(c *mongo.Collection) {
	deploymentCollection = c
	log.Println("Deployment collection initialized")
}

func CreateDeployment(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateDeployment called") 
	if deploymentCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}

	appName := r.FormValue("appName")
	userName := r.FormValue("userName")
	deployPath := r.FormValue("deployPath")
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if path exists
	if _, err := os.Stat(deployPath); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("Path does not exist: %s", deployPath), http.StatusBadRequest)
		return
	}

	// Save file
	filePath := filepath.Join(deployPath, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	// Save to DB
	deployment := models.Deployment{
		ID:         primitive.NewObjectID(),
		AppName:    appName,
		UserName:   userName,
		DeployPath: deployPath,
		FileName:   header.Filename,
		CreatedAt:  time.Now(),
	}

	_, err = deploymentCollection.InsertOne(context.Background(), deployment)
	if err != nil {
		http.Error(w, "Failed to save deployment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deployment)
}

func GetAllDeployments(w http.ResponseWriter, r *http.Request) {
	cursor, err := deploymentCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch deployments", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var deployments []models.Deployment
	if err := cursor.All(context.Background(), &deployments); err != nil {
		http.Error(w, "Failed to decode deployments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deployments)
}

func GetDeploymentByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/deployments/"):]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid deployment ID", http.StatusBadRequest)
		return
	}

	var deployment models.Deployment
	err = deploymentCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&deployment)
	if err != nil {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(deployment)
}
