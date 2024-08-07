# Overview of How to Use the APIs 

Use the following command to start the server. The server runs on the localhost:8080
``` bash
go run main.go
```

* GET Requests
  - http://localhost:8080/tasks : Get tasks
  - http://localhost:8080/tasks/ID : Get task with taskId ID

* PUT Request
  - http://localhost:8080/tasks/ID: Update the fields of task with taskId ID

* DELETE Request
  - http://localhost:8080/tasks/ID: Delete the task with taskId ID

* POST Request
  - http://localhost:8080/tasks: Add new task


### This is the [API Documentation](https://documenter.getpostman.com/view/37363410/2sA3rzJCHr) for this Simple Task Mangement System.

