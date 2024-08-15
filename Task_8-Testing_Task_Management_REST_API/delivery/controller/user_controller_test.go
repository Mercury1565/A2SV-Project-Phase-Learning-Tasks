package controller

import (
	"Task_8-Testing_Task_Management_REST_API/bootstrap"
	"Task_8-Testing_Task_Management_REST_API/domain"
	"Task_8-Testing_Task_Management_REST_API/infrastructure"
	"Task_8-Testing_Task_Management_REST_API/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserControllerTestSuite struct {
	suite.Suite
	mockUserUsecase *mocks.UserUsecase
	controller      *UserController
	router          *gin.Engine
}

func (suite *UserControllerTestSuite) SetupSuite() {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		suite.Fail("Failed to load .env.test file", err)
	}

	suite.mockUserUsecase = new(mocks.UserUsecase)
	suite.controller = &UserController{
		UserUsecase: suite.mockUserUsecase,
		Env:         bootstrap.NewEnv(),
	}
	suite.router = gin.Default()

	// define the routes
	suite.router.POST("/register", suite.controller.HandelUserRegister)
	suite.router.POST("/login", suite.controller.HandelUserLogin)
	suite.router.PUT("/promote/:id", suite.controller.HandleUserPromotion)
}

func (suite *UserControllerTestSuite) TearDownTest() {
	suite.mockUserUsecase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestHandelUserRegister_Success() {
	requestUser := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "USER",
	}

	suite.mockUserUsecase.On("AreThereAnyUsers", mock.Anything).Return(false, nil).Once()
	suite.mockUserUsecase.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Once()
	suite.mockUserUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

	jsonUser, _ := json.Marshal(requestUser)
	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonUser))
	request.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusOK, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "user registered successfully")
}

func (suite *UserControllerTestSuite) TestHandelUserRegister_UserAlreadyExists() {
	requestUser := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "USER",
	}

	suite.mockUserUsecase.On("AreThereAnyUsers", mock.Anything).Return(false, nil).Once()
	suite.mockUserUsecase.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(&domain.User{}, nil).Once()

	jsonUser, _ := json.Marshal(requestUser)
	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonUser))
	request.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusNotAcceptable, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "user already exists")
}

func (suite *UserControllerTestSuite) TestHandleUserLogin_Success() {
	mockUser := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "USER",
	}

	hashedPassword, _ := infrastructure.HashPassword(mockUser.Password)
	mockUser.Password = string(hashedPassword)

	suite.mockUserUsecase.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
	suite.mockUserUsecase.On("CreateAccessToken", mock.AnythingOfType("*domain.User"), mock.Anything, mock.Anything).Return("mocked_jwt_token", nil).Once()

	requestUser := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonUser, _ := json.Marshal(requestUser)
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonUser))
	request.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusOK, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "user logged in successfully")
	suite.Contains(responseWriter.Body.String(), "mocked_jwt_token")
}

func (suite *UserControllerTestSuite) TestHandleUserLogin_UserNonExistent() {

	suite.mockUserUsecase.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Once()

	requestUser := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonUser, _ := json.Marshal(requestUser)
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonUser))
	request.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusNotFound, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "user doesn't exist")
}

func (suite *UserControllerTestSuite) TestHandleUserLogin_WrongPassword() {
	mockUser := &domain.User{
		Email:    "test@example.com",
		Password: "correct_password",
		Name:     "Test User",
		Role:     "USER",
	}

	hashedPassword, _ := infrastructure.HashPassword(mockUser.Password)
	mockUser.Password = string(hashedPassword)

	suite.mockUserUsecase.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()

	requestUser := &domain.User{
		Email:    "test@example.com",
		Password: "wrong password",
	}

	jsonUser, _ := json.Marshal(requestUser)
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonUser))
	request.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusUnauthorized, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "Incorrect password")
}

func (suite *UserControllerTestSuite) TestHandleUserPromotion_Success() {
	mockUser := &domain.User{
		UserID: primitive.NewObjectID(),
		Email:  "test@example.com",
		Role:   "USER",
	}

	suite.mockUserUsecase.On("GetByID", mock.Anything, mockUser.UserID.Hex()).Return(mockUser, nil).Once()
	suite.mockUserUsecase.On("UpdateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

	request, _ := http.NewRequest(http.MethodPut, "/promote/"+mockUser.UserID.Hex(), nil)
	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusOK, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "user promoted to admin status")
}

func (suite *UserControllerTestSuite) TestHandleUserPromotion_UserNonExistent() {

	suite.mockUserUsecase.On("GetByID", mock.Anything, "nonExistingID").Return(nil, nil).Once()

	request, _ := http.NewRequest(http.MethodPut, "/promote/nonExistingID", nil)
	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusNotFound, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "user not found")
}

func (suite *UserControllerTestSuite) TestHandleUserPromotion_UserAlreadyAdmin() {
	mockUser := &domain.User{
		UserID: primitive.NewObjectID(),
		Email:  "test@example.com",
		Role:   "ADMIN",
	}

	suite.mockUserUsecase.On("GetByID", mock.Anything, mockUser.UserID.Hex()).Return(mockUser, nil).Once()

	request, _ := http.NewRequest(http.MethodPut, "/promote/"+mockUser.UserID.Hex(), nil)
	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusOK, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "user is already an admin")
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
