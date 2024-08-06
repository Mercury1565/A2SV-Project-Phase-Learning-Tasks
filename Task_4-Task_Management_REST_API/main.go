package main

import "Task_4-Task_Management_REST_API/router"

func main() {
	route := router.SetUpRouter()
	route.Run("localhost:8080")
}
