package controller

import (
	"Task_8-Testing_Task_Management_REST_API/bootstrap"
	"Task_8-Testing_Task_Management_REST_API/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
	Env         *bootstrap.Env
}

func (controller *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := controller.TaskUsecase.GetTasks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (controller *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := controller.TaskUsecase.GetTaskByID(c, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (controller *TaskController) CreateTask(c *gin.Context) {
	var new_task domain.Task

	if e := c.BindJSON(&new_task); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := controller.TaskUsecase.Create(c, &new_task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task added successfully"})
}

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

func (controller *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := controller.TaskUsecase.DeleteTask(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
