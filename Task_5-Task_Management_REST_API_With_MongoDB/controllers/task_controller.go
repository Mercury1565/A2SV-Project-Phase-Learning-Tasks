package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Task_5-Task_Management_REST_API_With_MongoDB/data"
	"Task_5-Task_Management_REST_API_With_MongoDB/models"
)

// declare a taskManager instance
var taskManager *data.TaskManager

func InitializeTaskManger() {
	// initialize a taskManager instance
	taskManager = data.NewTaskManager()
}

// GetAllTasks returns all tasks.
func GetAllTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, taskManager.GetAllTasks())
}

// GetTask retrieves a task by its ID from the task collection and returns it as JSON.
// If the task is not found, it returns a JSON response with a "task not found" message.
func GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := taskManager.GetTask(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

// UpdateTask updates a task with the given ID.
// It expects a JSON request body containing the updated task details.
// If the request body is invalid, it returns a 400 Bad Request error with a "invalid request body" error message.
// If the task is not found, it returns a 404 Not Found error with a "task not found" message.
// Otherwise, it returns a 200 OK response with a success message.
func UpdateTask(c *gin.Context) {
	var updated_task models.AddedTask
	id := c.Param("id")

	if e := c.BindJSON(&updated_task); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := taskManager.UpdateTask(id, &updated_task)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// DeleteTask deletes a task with the given ID.
// If the task is not found, it returns a 404 Not Found error with a "task not found" message.
// Otherwise, it returns a 200 OK response with a success message.
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := taskManager.DeleteTask(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

// CreateTask creates a new task based on the JSON data provided in the request body.
// It binds the JSON data to the `new_task` variable and adds it to the task collection.
// If the request body is invalid, it returns a JSON response with an error message.
// If the task is added successfully, it returns a JSON response with a success message.
func CreateTask(c *gin.Context) {
	var new_task models.AddedTask

	if e := c.BindJSON(&new_task); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := taskManager.AddTask(&new_task)

	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "task added successfully"})
}
