package infrastructure

import (
	"Task_8-Testing_Task_Management_REST_API/bootstrap"
	"Task_8-Testing_Task_Management_REST_API/domain"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareSuite struct {
	suite.Suite
	router   *gin.Engine
	secret   string
	mockUser *domain.User
}

func (suite *AuthMiddlewareSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.router = gin.Default()
	suite.secret = bootstrap.NewEnv(2).AccessTokenSecret
	suite.mockUser = &domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "USER",
	}
}

func (suite *AuthMiddlewareSuite) TestJWTAuthMiddleware_Success() {
	suite.router.Use(JWTAuthMiddleware(suite.secret))
	suite.router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	accessToken, err := CreateAccessToken(suite.mockUser, suite.secret, bootstrap.NewEnv(2).AccessTokenExpiryHour)
	if err != nil {
		suite.Fail("Failed to generate token", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusOK, response.Code)
}

func (suite *AuthMiddlewareSuite) TestJWTAuthMiddleware_NoAuthHeader() {
	suite.router.Use(JWTAuthMiddleware(suite.secret))
	suite.router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusUnauthorized)
	})

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)

	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Contains(response.Body.String(), "authorization header is required")
}

func (suite *AuthMiddlewareSuite) TestJWTAuthMiddleware_InvalidAuthHeader() {
	suite.router.Use(JWTAuthMiddleware(suite.secret))
	suite.router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusUnauthorized)
	})

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", "invalid token")

	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Contains(response.Body.String(), "invalid authorization header")
}

func (suite *AuthMiddlewareSuite) TestJWTAuthMiddleware_UnauthorizedToken() {
	suite.router.Use(JWTAuthMiddleware(suite.secret))
	suite.router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusUnauthorized)
	})

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJpc3JhdEBnbWFpbC5jb20iLCJleHBpcm"
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Contains(response.Body.String(), "unauthorized user")
}

func (suite *AuthMiddlewareSuite) TestJWTAuthMiddleware_TokenExpired() {
	suite.router.Use(JWTAuthMiddleware(suite.secret))
	suite.router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusUnauthorized)
	})

	accessToken, err := CreateAccessToken(suite.mockUser, suite.secret, -1)
	if err != nil {
		suite.Fail("Failed to generate token", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	response := httptest.NewRecorder()
	suite.router.ServeHTTP(response, request)

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Contains(response.Body.String(), "token has expired")
}

func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareSuite))
}
