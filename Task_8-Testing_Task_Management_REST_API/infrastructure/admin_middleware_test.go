package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AdminMiddlewareSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *AdminMiddlewareSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
	suite.router.Use(AuthenticateAdmin())
}
func (suite *AdminMiddlewareSuite) TestAdminMiddleware_AdminUser() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("claims", jwt.MapClaims{"role": "ADMIN"})

	AuthenticateAdmin()(c)

	suite.Equal(http.StatusOK, c.Writer.Status())
}

func (suite *AdminMiddlewareSuite) TestAdminMiddleware_NonAdminUser() {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Set("claims", jwt.MapClaims{"role": "USER"})

	AuthenticateAdmin()(c)

	suite.Equal(http.StatusUnauthorized, c.Writer.Status())
	suite.Contains(recorder.Body.String(), "unauthorized user")
}

func TestAdminMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AdminMiddlewareSuite))
}
