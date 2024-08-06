package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"Task_4-Task_Management_REST_API/data"
	"Task_4-Task_Management_REST_API/models"
)

// GetAllTasks returns all tasks.
func GetAllTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data.Tasks)
}

// GetTaskById retrieves a task from the data.Tasks slice based on the provided ID.
// It returns a pointer to the task and an error if the task is not found.
func GetTaskById(id string) (*models.Task, error) {
	for idx, task := range data.Tasks {
		if task.ID == id {
			return &data.Tasks[idx], nil
		}
	}
	return nil, errors.New("task not found")
}

// GetTask retrieves a task by its ID and returns it as JSON.
// If the task is not found, it returns a JSON response with a "Task not found" message.
func GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := GetTaskById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
	}
	c.IndentedJSON(http.StatusOK, task)
}

// UpdateTask updates a task with the provided ID.
// It retrieves the task by ID, then updates the task's properties based on the provided JSON payload.
// The updated task is then returned as a JSON response.
func UpdateTask(c *gin.Context) {
	var updated_task models.Task
	id := c.Param("id")

	task, err := GetTaskById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	e := c.BindJSON(&updated_task)

	if e != nil {
		return
	}

	// update the task
	if updated_task.Title != "" {
		task.Title = updated_task.Title
	}
	if updated_task.Description != "" {
		task.Description = updated_task.Description
	}
	if updated_task.DueDate != "" {
		task.DueDate = updated_task.DueDate
	}
	if updated_task.Status != "" {
		task.Status = updated_task.Status
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// DeleteTask deletes a task with the given ID from the data.Tasks slice.
// If the task is found and deleted, it returns the updated list of tasks.
// If the task is not found, it returns a JSON response with a "Task not found" message.
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	for i, task := range data.Tasks {
		if id == task.ID {
			data.Tasks = append(data.Tasks[:i], data.Tasks[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// CreateTask creates a new task based on the JSON data provided in the request body.
// It binds the JSON data to the `newTask` variable, appends it to the `data.Tasks` slice,
// and returns the updated list of tasks as a JSON response.
func CreateTask(c *gin.Context) {
	var newTask models.Task

	err := c.BindJSON(&newTask)

	if err != nil {
		return
	}

	data.Tasks = append(data.Tasks, newTask)
	c.IndentedJSON(http.StatusOK, data.Tasks)
}
