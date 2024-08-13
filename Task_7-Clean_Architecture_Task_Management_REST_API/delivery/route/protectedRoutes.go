package route

import (
	"Task_7-Clean_Architecture_Task_Management_REST_API/bootstrap"
	"Task_7-Clean_Architecture_Task_Management_REST_API/delivery/controller"
	"Task_7-Clean_Architecture_Task_Management_REST_API/domain"
	"Task_7-Clean_Architecture_Task_Management_REST_API/repository"
	"Task_7-Clean_Architecture_Task_Management_REST_API/usecases"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewProtectedRouter(env *bootstrap.Env, timeout time.Duration, database mongo.Database, group *gin.RouterGroup) {
	taskRepo := repository.NewTaskRepo(database, domain.CollectionTask)

	protectedRouteTaskController := &controller.TaskController{
		TaskUsecase: usecases.NewTaskUsecase(taskRepo, timeout),
		Env:         env,
	}

	group.GET("/tasks", protectedRouteTaskController.GetAllTasks)
	group.GET("/tasks/:id", protectedRouteTaskController.GetTask)
}
