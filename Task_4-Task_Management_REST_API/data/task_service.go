package data

import (
	"Task_4-Task_Management_REST_API/models"
	"time"
)

// define in-memory dummy data
var Tasks = []models.Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now().Format("2006-01-02"), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2).Format("2006-01-02"), Status: "Completed"},
}
