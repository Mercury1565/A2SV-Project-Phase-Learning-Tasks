package usecases

import (
	"Task_8-Testing_Task_Management_REST_API/domain"
	"Task_8-Testing_Task_Management_REST_API/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUsecaseTestSuite struct {
	suite.Suite
	taskUsecase  *taskUsecase
	taskMockRepo *mocks.TaskRepository
}

// setupSuite runs once before all tests in the suite
func (suite *TaskUsecaseTestSuite) SetupSuite() {
	suite.taskMockRepo = new(mocks.TaskRepository)
	suite.taskUsecase = &taskUsecase{
		taskRepository: suite.taskMockRepo,
		contextTimeout: time.Second * 2,
	}
}

func (suite *TaskUsecaseTestSuite) TearDownSuite() {
	suite.taskMockRepo.AssertExpectations(suite.T())
}

func (suite *TaskUsecaseTestSuite) TestCreate() {
	mockTask := &domain.Task{
		Title:       "test title",
		Description: "test description",
		DueDate:     time.Now().UTC().Truncate(time.Second),
		Status:      "test status",
	}

	suite.taskMockRepo.On("Create", mock.Anything, mockTask).Return(nil)

	err := suite.taskUsecase.Create(context.Background(), mockTask)

	assert.NoError(suite.T(), err)
}

func (suite *TaskUsecaseTestSuite) TestGetTasks() {
	mockTasks := []domain.Task{
		{
			ID:          primitive.NewObjectID(),
			Title:       "test title 1",
			Description: "test description 1",
			DueDate:     time.Now().UTC().Truncate(time.Second),
			Status:      "test status 1",
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "test title 2",
			Description: "test description 2",
			DueDate:     time.Now().UTC().Truncate(time.Second),
			Status:      "test status 2",
		},
	}

	suite.taskMockRepo.On("GetTasks", mock.Anything).Return(mockTasks, nil)
	tasks, err := suite.taskUsecase.GetTasks(context.Background())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockTasks, tasks)
}

func (suite *TaskUsecaseTestSuite) TestGetTaskByID() {
	mockTask := domain.Task{
		Title:       "test title",
		Description: "test description",
		DueDate:     time.Now().UTC().Truncate(time.Second),
		Status:      "test status",
	}

	suite.taskMockRepo.On("GetTaskByID", mock.Anything, mockTask.ID.Hex()).Return(mockTask, nil)

	task, err := suite.taskUsecase.GetTaskByID(context.Background(), mockTask.ID.Hex())

	// assert no error occured
	assert.NoError(suite.T(), err)

	// assert 'mockTask' is equal to 'user'
	assert.Equal(suite.T(), mockTask, task)
}

func (suite *TaskUsecaseTestSuite) TestUpdateTask() {
	mockTask := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "test title",
		Description: "test description",
		DueDate:     time.Now().UTC().Truncate(time.Second),
		Status:      "test status",
	}

	suite.taskMockRepo.On("UpdateTask", mock.Anything, mockTask.ID.Hex(), mockTask).Return(nil)

	err := suite.taskUsecase.UpdateTask(context.Background(), mockTask.ID.Hex(), mockTask)

	// assert no error occured
	assert.NoError(suite.T(), err)
}

func (suite *TaskUsecaseTestSuite) TestDeleteTask() {
	mockTask := &domain.Task{
		ID: primitive.NewObjectID(),
	}

	suite.taskMockRepo.On("DeleteTask", mock.Anything, mockTask.ID.Hex()).Return(nil)

	err := suite.taskUsecase.DeleteTask(context.Background(), mockTask.ID.Hex())

	// assert no error occured
	assert.NoError(suite.T(), err)
}

// TestUserUsecaseTestSuite runs the test suite
func TestTaskUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}
