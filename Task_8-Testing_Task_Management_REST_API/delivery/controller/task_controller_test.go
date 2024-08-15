package controller

import (
	"Task_8-Testing_Task_Management_REST_API/domain"
	"Task_8-Testing_Task_Management_REST_API/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskControllerTestSuite struct {
	suite.Suite
	mockTaskUsecase *mocks.TaskUsecase
	controller      *TaskController
	router          *gin.Engine
}

func (suite *TaskControllerTestSuite) SetupSuite() {
	suite.mockTaskUsecase = new(mocks.TaskUsecase)
	suite.controller = &TaskController{
		TaskUsecase: suite.mockTaskUsecase,
	}
	suite.router = gin.Default()

	// define the rotes
	suite.router.GET("/tasks", suite.controller.GetAllTasks)
	suite.router.GET("/tasks/:id", suite.controller.GetTask)
	suite.router.POST("/tasks", suite.controller.CreateTask)
	suite.router.PUT("/tasks/:id", suite.controller.UpdateTask)
	suite.router.DELETE("/tasks/:id", suite.controller.DeleteTask)
}

func (suite *TaskControllerTestSuite) TearDownSuite() {
	suite.mockTaskUsecase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestGetAllTasks_Success() {
	mockTasks := []domain.Task{
		{ID: primitive.NewObjectID(), Title: "Task 1", Description: "Task 1 Description", DueDate: time.Now(), Status: "test status"},
		{ID: primitive.NewObjectID(), Title: "Task 2", Description: "Task 2 Description", DueDate: time.Now(), Status: "test status"},
	}

	suite.mockTaskUsecase.On("GetTasks", mock.Anything).Return(mockTasks, nil).Once()

	request, _ := http.NewRequest(http.MethodGet, "/tasks", nil) // create a HTTP request to be passed to the handler
	responseWriter := httptest.NewRecorder()                     // declare a new HTTP response writer to be used later
	suite.router.ServeHTTP(responseWriter, request)              // make the HTTP request, HTTP response written in 'responseWriter'

	suite.Equal(http.StatusOK, responseWriter.Code)
}

func (suite *TaskControllerTestSuite) TestGetAllTasks_InternalServerError() {
	suite.mockTaskUsecase.On("GetTasks", mock.Anything).Return([]domain.Task{}, errors.New("internal server error")).Once()

	request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusInternalServerError, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "internal server error")
}

func (suite *TaskControllerTestSuite) TestGetTask_Success() {
	mockTask := domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "Test Task",
		Description: "Test Task Description",
		DueDate:     time.Now(),
		Status:      "Test Status",
	}

	suite.mockTaskUsecase.On("GetTaskByID", mock.Anything, mock.AnythingOfType("string")).Return(mockTask, nil).Once()

	request, _ := http.NewRequest(http.MethodGet, "/tasks/"+mockTask.ID.Hex(), nil)
	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusOK, responseWriter.Code)
}

func (suite *TaskControllerTestSuite) TestGetTask_NotFound() {
	suite.mockTaskUsecase.On("GetTaskByID", mock.Anything, mock.AnythingOfType("string")).Return(domain.Task{}, errors.New("task not found")).Once()

	request, _ := http.NewRequest(http.MethodGet, "/tasks/THIS_ID_DOESNT EXIST", nil)
	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusNotFound, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "task not found")
}

func (suite *TaskControllerTestSuite) TestCreateTask_Success() {
	mockTask := domain.Task{
		Title:       "Test Task",
		Description: "Test Task Description",
		DueDate:     time.Now(),
		Status:      "Test Status",
	}

	suite.mockTaskUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil).Once()

	jsonTask, _ := json.Marshal(mockTask) // marshal mockTask to JSON
	request, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonTask))
	request.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusOK, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "task added successfully")
}

func (suite *TaskControllerTestSuite) TestUpdateTask_Success() {
	updatedTask := domain.Task{
		Title:       "Updated Task",
		Description: "Updated Task Description",
		DueDate:     time.Now(),
		Status:      "Updated Status",
	}

	suite.mockTaskUsecase.On("UpdateTask", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Task")).Return(nil).Once()

	jsonTask, _ := json.Marshal(updatedTask)
	request, _ := http.NewRequest(http.MethodPut, "/tasks/"+updatedTask.ID.Hex(), bytes.NewBuffer(jsonTask))
	request.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusOK, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "task updated successfully")
}

func (suite *TaskControllerTestSuite) TestUpdateTask_TaskNotFound() {
	updatedTask := domain.Task{
		Title:       "Updated Task",
		Description: "Updated Task Description",
		DueDate:     time.Now(),
		Status:      "Updated Status",
	}

	suite.mockTaskUsecase.On("UpdateTask", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Task")).Return(errors.New("task not found")).Once()

	jsonTask, _ := json.Marshal(updatedTask)
	request, _ := http.NewRequest(http.MethodPut, "/tasks/"+updatedTask.ID.Hex(), bytes.NewBuffer(jsonTask))
	request.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusNotFound, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "task not found")
}

func (suite *TaskControllerTestSuite) TestDeleteTask_Success() {
	suite.mockTaskUsecase.On("DeleteTask", mock.Anything, mock.Anything).Return(nil).Once()

	request, _ := http.NewRequest(http.MethodDelete, "/tasks/JIBBER_JABBER", nil)
	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusOK, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "task deleted successfully")
}

func (suite *TaskControllerTestSuite) TestDeleteTask_TaskNotFound() {
	suite.mockTaskUsecase.On("DeleteTask", mock.Anything, mock.Anything).Return(errors.New("task not found")).Once()

	request, _ := http.NewRequest(http.MethodDelete, "/tasks/JIBBER_JABBER", nil)
	responseWriter := httptest.NewRecorder()
	suite.router.ServeHTTP(responseWriter, request)

	suite.Equal(http.StatusNotFound, responseWriter.Code)
	suite.Contains(responseWriter.Body.String(), "task not found")
}

func TestTaskControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
}
