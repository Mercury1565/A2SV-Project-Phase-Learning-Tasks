package data

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// instantiate database client
var client *mongo.Client

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
