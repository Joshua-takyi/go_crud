package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
) // Corrected import path for MongoDB driver

type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty"       bson:"_id"`
	Title       string             `json:"title"              bson:"title"       binding:"required,min=1,max=100"`
	Description string             `json:"description"        bson:"description" binding:"max=1000"`
	Image       []string           `json:"image"              bson:"image"`
	Priority    Priority           `json:"priority"           bson:"priority"    binding:"required"`
	Tags        []string           `json:"tags,omitempty"     bson:"tags"        binding:"dive,max=20"`
	Completed   bool               `json:"completed"          bson:"completed"`
	Metadata    Metadata           `json:"metadata"           bson:"metadata"`
}

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Metadata struct {
	CreatedAt time.Time `json:"created_at"             bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at"             bson:"updated_at"`
}
