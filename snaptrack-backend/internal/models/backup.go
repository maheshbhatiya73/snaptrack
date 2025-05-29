package models

import  (
	"go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Schedule struct {
	Kind string    `json:"kind" bson:"kind"`
	Date time.Time `json:"date" bson:"date"`
}

type Backup struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	App             string             `json:"app" bson:"app"`
	Type            string             `json:"type" bson:"type"`
	Size            string             `json:"size" bson:"size"`
	Status          string             `json:"status" bson:"status"`
	SourcePath      string             `json:"sourcePath" bson:"sourcePath"`
	DestinationPath string             `json:"destinationPath" bson:"destinationPath"`
	FileType        string             `json:"fileType" bson:"fileType"`
	Schedule        Schedule           `json:"schedule" bson:"schedule"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
}
