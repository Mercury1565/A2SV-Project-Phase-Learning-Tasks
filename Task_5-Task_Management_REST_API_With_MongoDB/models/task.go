package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a task in the task management system.
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"duedate" bson:"duedate"`
	Status      string             `json:"status" bson:"status"`
}

type AddedTask struct {
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	DueDate     time.Time `json:"duedate" bson:"duedate"`
	Status      string    `json:"status" bson:"status"`
}
