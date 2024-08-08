package main

import (
	"Task_5-Task_Management_REST_API_With_MongoDB/controllers"
	"Task_5-Task_Management_REST_API_With_MongoDB/data"
	"Task_5-Task_Management_REST_API_With_MongoDB/router"
)

// define essential parameters for MongoDB connection
const URI = "mongodb://localhost:27017"
const DATABASE_NAME = "test"
const COLLECTION_NAME = "tasks"

// define application URL
const URL = "localhost:8080"

func main() {
	data.StartMongoDB(URI, DATABASE_NAME, COLLECTION_NAME)
	defer data.CloseMongoDB()

	controllers.InitializeTaskManger()

	route := router.SetUpRouter()
	route.Run(URL)
}
