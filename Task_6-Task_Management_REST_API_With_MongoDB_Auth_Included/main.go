package main

import (
	"Task_6-Task_Management_REST_API_With_MongoDB/controllers"
	"Task_6-Task_Management_REST_API_With_MongoDB/data"
	"Task_6-Task_Management_REST_API_With_MongoDB/router"
)

// define essential parameters for MongoDB connection
const URI = "mongodb://localhost:27017"
const DATABASE_NAME = "test"
const TASK_COLLECTION_NAME = "tasks"
const USER_COLLECTION_NAME = "users"

// define application URL
const URL = "localhost:8080"

func main() {
	data.StartMongoDB(URI, DATABASE_NAME, TASK_COLLECTION_NAME, USER_COLLECTION_NAME)
	defer data.CloseMongoDB()

	controllers.InitializeSystemMangement()

	route := router.SetUpRouter()
	route.Run(URL)
}
