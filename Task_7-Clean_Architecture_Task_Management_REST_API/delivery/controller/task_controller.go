package controller

import (
	"Task_7-Clean_Architecture_Task_Management_REST_API/bootstrap"
	"Task_7-Clean_Architecture_Task_Management_REST_API/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
	Env         *bootstrap.Env
}

// GetAllTasks retrieves all tasks from the database and returns them as a JSON response.
func (controller *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := controller.TaskUsecase.GetTasks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask retrieves a task by its ID.
// It takes a gin.Context object and the task ID as parameters.
// It returns the retrieved task or an error if the task is not found.
func (controller *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := controller.TaskUsecase.GetTaskByID(c, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask is a method of the TaskController struct that handles the creation of a new task.
// It takes a gin.Context object as a parameter, which represents the HTTP request and response.
// The function first binds the JSON data from the request body to a new_task variable.
// If the request body is invalid, it returns a JSON response with an error message.
// Otherwise, it calls the Create method of the TaskUsecase to create the task.
// If an error occurs during the creation process, it returns a JSON response with the error message.
// Finally, it returns a JSON response with a success message if the task is created successfully.
func (controller *TaskController) CreateTask(c *gin.Context) {
	var new_task domain.Task

	if e := c.BindJSON(&new_task); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := controller.TaskUsecase.Create(c, &new_task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task added successfully"})
}

// UpdateTask updates a task with the given ID.
// It receives a JSON payload containing the updated task information.
// If the request body is invalid, it returns a 400 Bad Request response.
// If the task with the given ID is not found, it returns a 404 Not Found response.
// Otherwise, it updates the task and returns a 200 OK response.
func (controller *TaskController) UpdateTask(c *gin.Context) {
	var updated_task domain.Task
	id := c.Param("id")

	err := c.BindJSON(&updated_task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = controller.TaskUsecase.UpdateTask(c, id, &updated_task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// DeleteTask deletes a task with the given ID.
func (controller *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := controller.TaskUsecase.DeleteTask(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
