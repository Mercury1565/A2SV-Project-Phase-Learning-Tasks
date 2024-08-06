package models

// Task represents a task in the task management system.
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"duedate"`
	Status      string `json:"status"`
}
