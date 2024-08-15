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
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userUsecase  *userUsecase
	userMockRepo *mocks.UserRepository
}

// setupSuite runs once before all tests in the suite
func (suite *UserUsecaseTestSuite) SetupSuite() {
	suite.userMockRepo = new(mocks.UserRepository)
	suite.userUsecase = &userUsecase{
		userRepository: suite.userMockRepo,
		contextTimeout: time.Second * 2,
	}
}

func (suite *UserUsecaseTestSuite) TearDownSuite() {
	suite.userMockRepo.AssertExpectations(suite.T())
}

func (suite *UserUsecaseTestSuite) TestCreate() {
	mockUser := &domain.User{
		Name:     "test name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "test role",
	}

	suite.userMockRepo.On("Create", mock.Anything, mockUser).Return(nil)

	err := suite.userUsecase.Create(context.Background(), mockUser)

	assert.NoError(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestGetByEmail() {
	mockUser := &domain.User{
		Name:     "test name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "test role",
	}

	suite.userMockRepo.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil)

	user, err := suite.userUsecase.GetByEmail(context.Background(), mockUser.Email)

	// assert no error occured
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, user)
}

func (suite *UserUsecaseTestSuite) TestGetByID() {
	mockUser := &domain.User{
		Name:     "test name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "test role",
	}

	suite.userMockRepo.On("GetByID", mock.Anything, mockUser.UserID.Hex()).Return(mockUser, nil)

	user, err := suite.userUsecase.GetByID(context.Background(), mockUser.UserID.Hex())

	// assert no error occured
	assert.NoError(suite.T(), err)

	// assert 'mockUser' is equal to 'user'
	assert.Equal(suite.T(), mockUser, user)
}

func (suite *UserUsecaseTestSuite) TestUpdateUser() {
	mockUser := &domain.User{
		Name:     "test name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "test role",
	}

	suite.userMockRepo.On("UpdateUser", mock.Anything, mockUser).Return(nil)

	err := suite.userUsecase.UpdateUser(context.Background(), mockUser)

	// assert no error occured
	assert.NoError(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestAreThereAnyUsers() {
	// case 1: users exist
	suite.userMockRepo.On("AreThereAnyUsers", mock.Anything).Return(true, nil)
	result, err := suite.userUsecase.AreThereAnyUsers(context.Background())

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result)

	// case 2: no users exist
	suite.userMockRepo.On("AreThereAnyUsers", mock.Anything).Return(false, nil)
	result, err = suite.userUsecase.AreThereAnyUsers(context.Background())

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result)
}

func (suite *UserUsecaseTestSuite) TestCreateAccessToken() {
	mockUser := &domain.User{
		Name:     "test name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "test role",
	}

	token, err := suite.userUsecase.CreateAccessToken(mockUser, "secret", 3600)

	// assert no error occured
	assert.NoError(suite.T(), err)

	// assert token is not nil
	assert.NotEmpty(suite.T(), token)
}

// TestUserUsecaseTestSuite runs the test suite
func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
