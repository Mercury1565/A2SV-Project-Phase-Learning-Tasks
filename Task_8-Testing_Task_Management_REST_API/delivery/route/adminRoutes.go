package route

import (
	"Task_8-Testing_Task_Management_REST_API/bootstrap"
	"Task_8-Testing_Task_Management_REST_API/delivery/controller"
	"Task_8-Testing_Task_Management_REST_API/domain"
	"Task_8-Testing_Task_Management_REST_API/repository"
	"Task_8-Testing_Task_Management_REST_API/usecases"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAdminRouter(env *bootstrap.Env, timeout time.Duration, database mongo.Database, group *gin.RouterGroup) {
	userRepo := repository.NewUserRepo(database, domain.CollectionUser)
	taskRepo := repository.NewTaskRepo(database, domain.CollectionTask)

	adminRouteUserController := &controller.UserController{
		UserUsecase: usecases.NewUserUsecase(userRepo, timeout),
		Env:         env,
	}

	adminRouteTaskController := &controller.TaskController{
		TaskUsecase: usecases.NewTaskUsecase(taskRepo, timeout),
		Env:         env,
	}

	group.POST("/promote/:id", adminRouteUserController.HandleUserPromotion)
	group.POST("/tasks", adminRouteTaskController.CreateTask)
	group.PUT("/tasks/:id", adminRouteTaskController.UpdateTask)
	group.DELETE("/tasks/:id", adminRouteTaskController.DeleteTask)
}
