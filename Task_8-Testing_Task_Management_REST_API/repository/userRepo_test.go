package repository

import (
	"Task_8-Testing_Task_Management_REST_API/domain"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type UserRepoTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       *userRepo
	collection *mongo.Collection
}

// SetupSuite runs once before any test in the suite
func (suite *UserRepoTestSuite) SetupSuite() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		suite.T().Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		suite.T().Fatalf("Failed to ping MongoDB: %v", err)
	}

	suite.db = client.Database("test_db")
	suite.repo = &userRepo{
		database:   *suite.db,
		collection: "test_users",
	}
	suite.collection = suite.db.Collection("test_users")
}

// TearDownSuite runs once after all tests in the suite have finished
func (suite *UserRepoTestSuite) TearDownSuite() {
	if err := suite.db.Drop(context.Background()); err != nil {
		suite.T().Fatalf("Failed to drop test database: %v", err)
	}
	if err := suite.db.Client().Disconnect(context.Background()); err != nil {
		suite.T().Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
}

// setup tests before each test
func (suite *UserRepoTestSuite) SetupTest() {
	// clear the task collection before each test
	suite.collection.Drop(context.Background())
}

func (suite *UserRepoTestSuite) TestCreateUser() {
	user := &domain.User{
		Name:     "Test Name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "Test Role",
	}

	// check if user is created without error
	err := suite.repo.Create(context.Background(), user)
	suite.NoError(err)

	var insertedUser domain.User

	// check if user is indeed added to the collection
	err = suite.collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&insertedUser)
	suite.NoError(err)

	// check if user has the right fields
	suite.Equal(user.Name, insertedUser.Name)
	suite.Equal(user.Email, insertedUser.Email)
	suite.Equal(user.Password, insertedUser.Password)
	suite.Equal(user.Role, insertedUser.Role)
}

func (suite *UserRepoTestSuite) TestGetByEmail() {
	user := &domain.User{
		Name:     "Test Name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "Test Role",
	}

	// check if user is created without error
	err := suite.repo.Create(context.Background(), user)
	suite.NoError(err)

	// check if user is retireved without error and retrieved user is not NIL
	retrievedUser, err := suite.repo.GetByEmail(context.Background(), user.Email)
	suite.NoError(err)
	suite.NotNil(retrievedUser)

	// check if retrieved user has the right fields
	suite.Equal(user.Name, retrievedUser.Name)
	suite.Equal(user.Email, retrievedUser.Email)
	suite.Equal(user.Password, retrievedUser.Password)
	suite.Equal(user.Role, retrievedUser.Role)
}

func (suite *UserRepoTestSuite) TestGetByID() {
	user := &domain.User{
		Name:     "Test Name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "Test Role",
	}

	// check if user is created without error
	err := suite.repo.Create(context.Background(), user)
	suite.NoError(err)

	// check if user is retireved without error and retrieved user is not NIL
	retrievedUser, err := suite.repo.GetByID(context.Background(), user.UserID.Hex())
	suite.NoError(err)
	suite.NotNil(retrievedUser)

	// check if retrieved user has the right fields
	suite.Equal(user.Name, retrievedUser.Name)
	suite.Equal(user.Email, retrievedUser.Email)
	suite.Equal(user.Password, retrievedUser.Password)
	suite.Equal(user.Role, retrievedUser.Role)
}

func (suite *UserRepoTestSuite) TestUpdateUser() {
	user := &domain.User{
		Name:     "Test Name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "Test Role",
	}

	// check if user is created without error
	err := suite.repo.Create(context.Background(), user)
	suite.NoError(err)

	user.Name = "newName"
	user.Email = "newTest@example.com"
	user.Password = "newpassword123"
	user.Role = "newRole"

	// check if user is updated without error
	err = suite.repo.UpdateUser(context.Background(), user)
	suite.NoError(err)

	// check if user is found in the collection
	var updatedUser domain.User
	err = suite.collection.FindOne(context.Background(), bson.M{"_id": user.UserID}).Decode(&updatedUser)
	suite.NoError(err)

	// check if the fields of the updated user are correct
	suite.Equal("newName", updatedUser.Name)
	suite.Equal("newTest@example.com", updatedUser.Email)
	suite.Equal("newpassword123", updatedUser.Password)
	suite.Equal("newRole", updatedUser.Role)
}

func (suite *UserRepoTestSuite) TestAreThereAnyUsers() {
	// first check with no users
	checkUsers, err := suite.repo.AreThereAnyUsers(context.Background())
	suite.NoError(err)
	suite.False(checkUsers)

	// add new user
	user := &domain.User{
		Name:     "Test Name",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "Test Role",
	}

	// check if user is created without error
	err = suite.repo.Create(context.Background(), user)
	suite.NoError(err)

	// check after adding new user
	checkUsers, err = suite.repo.AreThereAnyUsers(context.Background())
	suite.NoError(err)
	suite.True(checkUsers)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
