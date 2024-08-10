package data

import (
	"Task_6-Task_Management_REST_API_With_MongoDB/models"
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

// Global variable to store the JWT secret
var JwtSecret = []byte("your_jwt_secret")

// define mongoDB task_collection for tasks
var task_collection *mongo.Collection

// define mongoDB task_collection for tasks
var user_collection *mongo.Collection

// instantiate database client
var client *mongo.Client

// define Task Manager Interface
type System_Manager_Interface interface {
	GetAllTasks()
	GetTask()
	UpdateTask()
	DeleteTask()
	CreateTask()
}

// define SystemManagement struct
type SystemManagement struct {
	task_collection *mongo.Collection
	user_collection *mongo.Collection
}

// define SystemManagement constructor
func NewSystemManager() *SystemManagement {
	return &SystemManagement{
		task_collection: task_collection,
		user_collection: user_collection,
	}
}

func StartMongoDB(URI string, DATABASE_NAME string, TASK_COLLECTION_NAME string, USER_COLLECTION_NAME string) {
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

	// instantiate the task_collection with name 'tasks' from the database 'test'
	task_collection = client.Database(DATABASE_NAME).Collection(TASK_COLLECTION_NAME)
	user_collection = client.Database(DATABASE_NAME).Collection(USER_COLLECTION_NAME)
	fmt.Println("Connected to MongoDB!")
}

func CloseMongoDB() {
	err := client.Disconnect(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func ConvertToObjectID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		// invalid id entered
		return primitive.NilObjectID, err
	}

	return objID, nil
}

// GetAllTasks retrieves all tasks from the task task_collection.
// It returns a slice of models.Task containing all the tasks.
func (taskManager *SystemManagement) GetAllTasks() []models.Task {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all tasks in task_collection
	cursor, err := taskManager.task_collection.Find(ctx, bson.M{})
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

// GetTask retrieves a task from the task task_collection based on the provided task ID.
// It returns the retrieved task and an error, if any.
func (taskManager *SystemManagement) GetTask(task_id string) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	obj_ID, err := ConvertToObjectID(task_id)
	var task models.Task

	if err != nil {
		return task, errors.New("invalid id entered")
	}

	err = taskManager.task_collection.FindOne(ctx, bson.M{"_id": obj_ID}).Decode(&task)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// task with id 'task_id' not found
			return task, fmt.Errorf("task with id '%v' not found", task_id)
		}
		log.Fatal(err)
	}

	return task, nil
}

// UpdateTask updates a task with the given task_id in the task task_collection.
// It takes the task_id string and the updated_task pointer to models.Task as parameters.
// It returns an error if the update operation fails.
func (taskManager *SystemManagement) UpdateTask(task_id string, updated_task *models.AddedTask) error {
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
	updateResult, err := taskManager.task_collection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Fatal(err)
	}

	if updateResult.MatchedCount == 0 {
		// task with id 'task_id' not found
		return fmt.Errorf("task with id '%v' not found", task_id)
	}

	return nil
}

// DeleteTask deletes a task from the task task_collection based on the given task ID.
// It returns an error if the task is not found or if there is an error during the deletion process.
func (taskManager *SystemManagement) DeleteTask(task_id string) error {
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
	deleteResult, err := taskManager.task_collection.DeleteOne(ctx, filter)

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

// AddTask adds a new task to the task task_collection.
// It takes a pointer to a models.Task struct as a parameter.
// It returns an error if there was a problem inserting the task into the task_collection.
func (taskManager *SystemManagement) AddTask(new_task *models.AddedTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// enforce user not to use empty task title
	if new_task.Title == "" {
		return errors.New("empty task title not allowed")
	}

	// insert new_task into taskManager
	_, err := taskManager.task_collection.InsertOne(ctx, new_task)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// AddNewUser adds a new user to the system.
// It takes a pointer to a `NewUser` struct as a parameter and returns an error.
// The function inserts the new user into the user collection in the database.
// If an error occurs during the insertion, it returns a `server error` error.
func (userManager *SystemManagement) AddNewUser(new_user *models.NewUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userManager.user_collection.InsertOne(ctx, new_user)

	if err != nil {
		return errors.New("server error")
	}

	return nil
}

// AreThereAnyUsers checks if there are any users in the system.
// It returns true if there are users, false otherwise.
// An error is returned if there was a problem checking for users.
func (userManager *SystemManagement) AreThereAnyUsers() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userManager.user_collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, nil
	}

	return count > 0, nil
}

// GetExistingUserByEmail retrieves an existing user from the database based on the provided email.
// If the user is found, it returns a pointer to the User object and nil error.
// If the user is not found, it returns nil and nil error.
// If any error occurs during the retrieval process, it returns nil and the error.
func (userManager *SystemManagement) GetExistingUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := userManager.user_collection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &existingUser, nil
}

// GetExistingUserByID retrieves an existing user from the database based on the provided ID.
// If the user is found, it returns a pointer to the User object and nil error.
// If the user is not found, it returns nil and nil error.
// If any other error occurs during the retrieval process, it returns nil and the error.
func (userManager *SystemManagement) GetExistingUserByID(id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := userManager.user_collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &existingUser, nil
}

// GetUser retrieves a user from the user collection based on the provided email.
// If the user is found, it returns an error indicating that the user already exists.
// If the user is not found, it returns nil.
func (userManager *SystemManagement) GetUser(new_user *models.NewUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userManager.user_collection.FindOne(ctx, bson.M{"email": new_user.Email}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
	}
	return fmt.Errorf("user with email %s already exists", new_user.Email)
}

// UpdateUser updates the information of a user in the system.
// It takes a pointer to a User struct as a parameter and returns an error if any.
func (userManager *SystemManagement) UpdateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.UserID}
	update := bson.M{"$set": user}

	_, err := userManager.user_collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
