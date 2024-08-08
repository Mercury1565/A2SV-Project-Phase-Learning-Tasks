package router

import (
	"Task_5-Task_Management_REST_API_With_MongoDB/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", controllers.GetAllTasks)
	r.GET("/tasks/:id", controllers.GetTask)
	r.PUT("/tasks/:id", controllers.UpdateTask)
	r.DELETE("/tasks/:id", controllers.DeleteTask)
	r.POST("/tasks", controllers.CreateTask)

	return r
}
