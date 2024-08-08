package data

import (
	"Task_5-Task_Management_REST_API_With_MongoDB/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// define mongoDB collection
var collection *mongo.Collection

// instantiate database client
var client *mongo.Client

// define Task Manager Interface
type Task_Manager_Interface interface {
	GetAllTasks()
	GetTask()
	UpdateTask()
	DeleteTask()
	CreateTask()
}

// define TaskManager struct
type TaskManager struct {
	collection *mongo.Collection
}

// define TaskManager constructor
func NewTaskManager() *TaskManager {
	return &TaskManager{
		collection: collection,
	}
}

func StartMongoDB(URI string, DATABASE_NAME string, COLLECTION_NAME string) {
	// set client options
	clientOptions := options.Client().ApplyURI(URI)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	// instantiate the collection with name 'tasks' from the database 'test'
	collection = client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
	fmt.Println("Connected to MongoDB!")
}

func CloseMongoDB() {
	err := client.Disconnect(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func ConvertToObjectID(task_id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(task_id)

	if err != nil {
		// invalid id entered
		return primitive.NilObjectID, err
	}

	return objID, nil
}

// GetAllTasks retrieves all tasks from the task collection.
// It returns a slice of models.Task containing all the tasks.
func (taskManager *TaskManager) GetAllTasks() []models.Task {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all tasks in collection
	cursor, err := taskManager.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var tasks []models.Task

	// decode all documents into 'tasks' slice
	err = cursor.All(ctx, &tasks)

	if err != nil {
		log.Fatal(err)
	}

	return tasks
}

// GetTask retrieves a task from the task collection based on the provided task ID.
// It returns the retrieved task and an error, if any.
func (taskManager *TaskManager) GetTask(task_id string) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	obj_ID, err := ConvertToObjectID(task_id)
	var task models.Task

	if err != nil {
		return task, errors.New("invalid id entered")
	}

	err = taskManager.collection.FindOne(ctx, bson.M{"_id": obj_ID}).Decode(&task)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// task with id 'task_id' not found
			return task, fmt.Errorf("task with id '%v' not found", task_id)
		}
		log.Fatal(err)
	}

	return task, nil
}

// UpdateTask updates a task with the given task_id in the task collection.
// It takes the task_id string and the updated_task pointer to models.Task as parameters.
// It returns an error if the update operation fails.
func (taskManager *TaskManager) UpdateTask(task_id string, updated_task *models.AddedTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	obj_ID, err := ConvertToObjectID(task_id)

	if err != nil {
		return errors.New("invalid id entered")
	}

	// initialize update parameter
	updated_fields := make(bson.M)

	// populate the update parameter by checking validity of the new task
	if updated_task.Title != "" {
		updated_fields["title"] = updated_task.Title
	}
	if updated_task.Description != "" {
		updated_fields["Description"] = updated_task.Description
	}
	if !updated_task.DueDate.IsZero() {
		updated_fields["duedate"] = updated_task.DueDate
	}
	if updated_task.Status != "" {
		updated_fields["status"] = updated_task.Status
	}

	// define update parameter
	update := bson.M{
		"$set": updated_fields,
	}

	// define filter parameter
	filter := bson.M{
		"_id": obj_ID,
	}

	// update the task with id 'task_id'
	updateResult, err := taskManager.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Fatal(err)
	}

	if updateResult.MatchedCount == 0 {
		// task with id 'task_id' not found
		return fmt.Errorf("task with id '%v' not found", task_id)
	}

	return nil
}

// DeleteTask deletes a task from the task collection based on the given task ID.
// It returns an error if the task is not found or if there is an error during the deletion process.
func (taskManager *TaskManager) DeleteTask(task_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	obj_ID, err := ConvertToObjectID(task_id)

	if err != nil {
		return errors.New("invalid id entered")
	}

	// define filter
	filter := bson.M{
		"_id": obj_ID,
	}

	// delete task with id 'task_id'
	deleteResult, err := taskManager.collection.DeleteOne(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}

	if deleteResult.DeletedCount == 0 {
		// task with id 'task_id' not found
		fmt.Println(task_id)
		return fmt.Errorf("task with id '%v' not found", task_id)
	}

	return nil
}

// AddTask adds a new task to the task collection.
// It takes a pointer to a models.Task struct as a parameter.
// It returns an error if there was a problem inserting the task into the collection.
func (taskManager *TaskManager) AddTask(new_task *models.AddedTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// enforce user not to use empty task title
	if new_task.Title == "" {
		return errors.New("empty task title not allowed")
	}

	// insert new_task into taskManager
	_, err := taskManager.collection.InsertOne(ctx, new_task)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
