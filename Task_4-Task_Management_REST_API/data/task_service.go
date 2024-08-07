package data

import (
	"Task_4-Task_Management_REST_API/models"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

// define Task Manager Interface
type TaskManager interface {
	GetAllTasks(c *gin.Context)
	GetTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	CreateTask(c *gin.Context)
}

// define task collection struct
type TaskCollection struct {
	Tasks []models.Task
}

// define TaskCollection constructor
func NewTaskCollection() *TaskCollection {
	return &TaskCollection{
		Tasks: []models.Task{
			{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now().Format("2006-01-02"), Status: "Pending"},
			{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"), Status: "In Progress"},
			{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2).Format("2006-01-02"), Status: "Completed"},
		},
	}
}

// GetTaskById retrieves a task from the data.Tasks slice based on the provided ID.
// It returns a pointer to the task and an error if the task is not found.
func (taskCollection *TaskCollection) GetTaskById(id string) (*models.Task, error) {
	for idx, task := range taskCollection.Tasks {
		if task.ID == id {
			return &taskCollection.Tasks[idx], nil
		}
	}
	return nil, errors.New("task not found")
}
