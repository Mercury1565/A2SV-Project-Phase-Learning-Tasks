package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionTask = "tasks"

type Task struct {
	ID          primitive.ObjectID `json:"-" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"duedate" bson:"duedate"`
	Status      string             `json:"status" bson:"status"`
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	GetTasks(c context.Context) ([]Task, error)
	GetTaskByID(c context.Context, taskID string) (Task, error)
	UpdateTask(c context.Context, taskID string, updated_task *Task) error
	DeleteTask(c context.Context, taskID string) error
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	GetTasks(c context.Context) ([]Task, error)
	GetTaskByID(c context.Context, taskID string) (Task, error)
	UpdateTask(c context.Context, taskID string, updated_task *Task) error
	DeleteTask(c context.Context, taskID string) error
}
