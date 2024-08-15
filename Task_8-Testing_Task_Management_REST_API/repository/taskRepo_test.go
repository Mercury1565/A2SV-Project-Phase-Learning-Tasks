package repository

import (
	"Task_8-Testing_Task_Management_REST_API/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type TaskRepoTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       *taskRepo
	collection *mongo.Collection
}

// SetupSuite runs once before any test in the suite
func (suite *TaskRepoTestSuite) SetupSuite() {
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
	suite.repo = &taskRepo{
		database:   *suite.db,
		collection: "test_tasks",
	}
	suite.collection = suite.db.Collection("test_tasks")
}

// TearDownSuite runs once after all tests in the suite have finished
func (suite *TaskRepoTestSuite) TearDownSuite() {
	if err := suite.db.Drop(context.Background()); err != nil {
		suite.T().Fatalf("Failed to drop test database: %v", err)
	}
	if err := suite.db.Client().Disconnect(context.Background()); err != nil {
		suite.T().Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
}

// setup tests before each test
func (suite *TaskRepoTestSuite) SetupTest() {
	// clear the task collection before each test
	suite.collection.Drop(context.Background())
}

func (suite *TaskRepoTestSuite) TestCreateTask() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     time.Now().UTC().Truncate(time.Second),
		Status:      "Test Status",
	}

	// check error during insertion
	err := suite.repo.Create(context.Background(), task)
	suite.NoError(err)

	// check if the task is indeed created and stored
	var insertedTask domain.Task
	err = suite.collection.FindOne(context.Background(), bson.M{"_id": task.ID}).Decode(&insertedTask)
	suite.NoError(err)

	// check if the created tasks contains the right parameters
	suite.Equal(task.Title, insertedTask.Title)
	suite.Equal(task.Description, insertedTask.Description)
	suite.Equal(task.DueDate, insertedTask.DueDate)
	suite.Equal(task.Status, insertedTask.Status)
}

func (suite *TaskRepoTestSuite) TestGetTasks() {
	tasks := []domain.Task{
		{Title: "Task 1", Description: "Description 1", DueDate: time.Now().UTC().Truncate(time.Second), Status: "Test Status 1"},
		{Title: "Task 2", Description: "Description 2", DueDate: time.Now().UTC().Truncate(time.Second), Status: "Test Status 2"},
	}

	// insert tasks in the collection
	for _, task := range tasks {
		err := suite.repo.Create(context.Background(), &task)
		suite.NoError(err)
	}

	retrievedTasks, err := suite.repo.GetTasks(context.Background())

	// check that tasks are retrieved without an error
	suite.NoError(err)

	// check if all the tasks are retrieved
	suite.Len(retrievedTasks, len(tasks))
}

func (suite *TaskRepoTestSuite) TestGetTaskByID() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     time.Now().UTC().Truncate(time.Second),
		Status:      "Test Status",
	}

	// check error during insertion
	err := suite.repo.Create(context.Background(), task)
	suite.NoError(err)

	retrievedTask, err := suite.repo.GetTaskByID(context.Background(), task.ID.Hex())

	// check that the task is retrieved without an error
	suite.NoError(err)

	// check if the retrieved task contains the right parameters
	suite.Equal(task.Title, retrievedTask.Title)
	suite.Equal(task.Description, retrievedTask.Description)
	suite.Equal(task.DueDate, retrievedTask.DueDate)
	suite.Equal(task.Status, retrievedTask.Status)
}

func (suite *TaskRepoTestSuite) TestUpdateTask() {
	originalTask := &domain.Task{
		Title:       "Original Task",
		Description: "Original Description",
		DueDate:     time.Now().UTC().Truncate(time.Second),
		Status:      "Original Status",
	}

	// check error during insertion
	err := suite.repo.Create(context.Background(), originalTask)
	suite.NoError(err)

	updatedTask := &domain.Task{
		Title:       "Updated Task",
		Description: "Updated Description",
		DueDate:     time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second),
		Status:      "Updated Status",
	}

	// check error during insertion
	err = suite.repo.UpdateTask(context.Background(), originalTask.ID.Hex(), updatedTask)
	suite.NoError(err)

	// check if the task is indeed created and stored
	var retrievedTask domain.Task
	err = suite.collection.FindOne(context.Background(), bson.M{"_id": originalTask.ID}).Decode(&retrievedTask)
	suite.NoError(err)

	// check if the updated tasks contains the right parameters
	suite.Equal(updatedTask.Title, retrievedTask.Title)
	suite.Equal(updatedTask.Description, retrievedTask.Description)
	suite.Equal(updatedTask.DueDate, retrievedTask.DueDate)
	suite.Equal(updatedTask.Status, retrievedTask.Status)
}

func (suite *TaskRepoTestSuite) TestDeleteTask() {
	task := &domain.Task{
		Title:       "Task to be Deleted",
		Description: "Description",
		DueDate:     time.Now().UTC().Truncate(time.Second),
		Status:      "Pending",
	}

	// check error during insertion
	err := suite.repo.Create(context.Background(), task)
	suite.NoError(err)

	// check task is deleted without error
	err = suite.repo.DeleteTask(context.Background(), task.ID.Hex())
	suite.NoError(err)

	// check task is removed from the collection
	err = suite.collection.FindOne(context.Background(), bson.M{"_id": task.ID}).Err()
	suite.Equal(mongo.ErrNoDocuments, err)
}

func TestTaskRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepoTestSuite))
}
