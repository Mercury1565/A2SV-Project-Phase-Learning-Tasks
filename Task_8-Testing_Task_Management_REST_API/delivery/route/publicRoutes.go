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

func NewPublicRouter(env *bootstrap.Env, timeout time.Duration, database mongo.Database, group *gin.RouterGroup) {
	userRepo := repository.NewUserRepo(database, domain.CollectionUser)

	publicRouteUserController := &controller.UserController{
		UserUsecase: usecases.NewUserUsecase(userRepo, timeout),
		Env:         env,
	}

	group.POST("/register", publicRouteUserController.HandelUserRegister)
	group.POST("/login", publicRouteUserController.HandelUserLogin)
}
