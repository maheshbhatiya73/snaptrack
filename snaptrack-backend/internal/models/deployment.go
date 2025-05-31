package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Deployment struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	AppName   string             `json:"appName" bson:"appName"`
	UserName  string             `json:"userName" bson:"userName"`
	DeployPath string            `json:"deployPath" bson:"deployPath"`
	FileName  string             `json:"fileName" bson:"fileName"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
