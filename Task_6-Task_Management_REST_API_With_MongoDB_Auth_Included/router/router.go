package router

import (
	"Task_6-Task_Management_REST_API_With_MongoDB/controllers"
	"Task_6-Task_Management_REST_API_With_MongoDB/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.HandelUserRegister)
	r.POST("/login", controllers.HandelUserLogin)

	// endpoints reserver for the ADMIN
	r.POST("/promote/:id", middleware.AuthenticateUser(), middleware.AuthenticateAdmin(), controllers.HandleUserPromotion)
	r.POST("/tasks", middleware.AuthenticateUser(), middleware.AuthenticateAdmin(), controllers.CreateTask)
	r.PUT("/tasks/:id", middleware.AuthenticateUser(), middleware.AuthenticateAdmin(), controllers.UpdateTask)
	r.DELETE("/tasks/:id", middleware.AuthenticateUser(), middleware.AuthenticateAdmin(), controllers.DeleteTask)

	// endpoints accessable by USER
	r.GET("/tasks", middleware.AuthenticateUser(), controllers.GetAllTasks)
	r.GET("/tasks/:id", middleware.AuthenticateUser(), controllers.GetTask)

	return r
}
