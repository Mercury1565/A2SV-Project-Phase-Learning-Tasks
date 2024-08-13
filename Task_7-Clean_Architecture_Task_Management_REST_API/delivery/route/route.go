package route

import (
	"Task_7-Clean_Architecture_Task_Management_REST_API/bootstrap"
	"Task_7-Clean_Architecture_Task_Management_REST_API/infrastructure"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	publicRouter := gin.Group("")
	protectedRouter := gin.Group("")
	adminRouter := gin.Group("")

	protectedRouter.Use(infrastructure.JWTAuthMiddleware(env.AccessTokenSecret))

	adminRouter.Use(
		infrastructure.JWTAuthMiddleware(env.AccessTokenSecret),
		infrastructure.AuthenticateAdmin(),
	)

	NewPublicRouter(env, timeout, db, publicRouter)
	NewProtectedRouter(env, timeout, db, protectedRouter)
	NewAdminRouter(env, timeout, db, adminRouter)
}
