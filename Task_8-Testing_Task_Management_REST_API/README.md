# Unit Testing the Task Managment System with GO | MongoDB | Authentication & Authorizaiton Features Included

## This repository shows how unit testing is performed on GO using the testify package. It is demonstrated using the simple Task Mangement App from Task-7 which was implemented with clean architecture.

## Prerequisites

- [Go](https://golang.org/doc/install)
- [MongoDB](https://docs.mongodb.com/manual/installation/)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/Mercury1565/A2SV-Project-Phase-Learning-Tasks
   ```

2. Navigate to the project directory:

   ```bash
   cd Task_8-Testing_Task_Management_REST_API
   ```

3. Download the required Go dependencies:

   ```bash
   go mod download
   ```

4. Tidy up the dependencies:
   ```bash
   go mod tidy
   ```

## MongoDB Setup

1. Start MongoDB:

   ```bash
   sudo systemctl start mongod
   ```

2. Verify that MongoDB has started successfully:

   ```bash
   sudo systemctl status mongod
   ```

3. Connect and use MongoDB:
   ```bash
   mongo
   ```

## Folder Structure

```
Task_manager/
├── bootstrap/
│   ├── app.go
│   ├── database.go
│   └── env.go
|
├── delivery/
│   ├── controllers/
│   │   ├── user_controller_test.go
│   │   ├── user_controller.go
│   │   └── tasks_controller_test.go
│   │   └── tasks_controller.go
│   ├── routers/
│   │   ├── adminRoutes.go
│   │   ├── protectedRoutes.go
│   │   ├── publicRoutes.go
│   │   └── route.go
│   └── main.go
|
├── domain/
│   ├── jwtCustom.go
│   ├── task.go
│   └── user.go
|
├── infrastructure/
│   ├── admin_middleware_test.go
│   ├── admin_middleware.go
│   ├── authenticate_middleware_test.go
│   ├── authenticate_middleware.go
│   ├── jwt_service.go
│   └── password_service.go
|
├── repository/
│   ├── taskRepo_test.go
│   ├── taskRepo.go
│   └── userRepo_test.go
│   └── userRepo.go
|
├── usecases/
│   ├── user_usecases_test.go
│   ├── user_usecases.go
│   └── task_usecases_test.go
│   └── task_usecases.go
|
├── README.md
├── go.mod
├── go.sum
└── .env

```

## Running the Project

1. Copy the `.env.example` file to a `.env` file:

   ```bash
   cp .env.example .env
   ```

2. Replace the placeholders in the `.env` file with your actual values.

3. Navigate to the delivery directory:

   ```bash
   cd delivery
   ```

4. Run the project:

   ```bash
   go run main.go
   ```

## API Endpoints

### APIs Related to Authentication

- POST Requests

  - http://localhost:8080/register: Register new user
  - http://localhost:8080/login : Authenticate and Signin Users
  - http://localhost:8080/promote/userID : Promote role of users to admin, only allowed for users with 'ADMIN' role

### APIs Related to task managment

- GET Requests

  - http://localhost:8080/tasks : Get tasks
  - http://localhost:8080/tasks/taskID : Get task with taskId ID

- PUT Request

  - http://localhost:8080/tasks/taskID: Update the fields of task with taskId ID, only allowed for users with 'ADMIN' role

- DELETE Request

  - http://localhost:8080/tasks/taskID: Delete the task with taskId ID, only allowed for users with 'ADMIN' role

- POST Request

  - http://localhost:8080/tasks: Add new task, only allowed for users with 'ADMIN' role

## Testing

Testing has been integrated into the project to ensure the reliability and correctness of the implemented functionalities.

Here's how you can run the tests

- Run all tests:

  ```bash
  go test -v ./...
  ```

- Run specific tests for the repository:

  ```bash
  go test ./repository -v
  ```

- Run specific tests for use cases:

  ```bash
  go test ./usecases -v
  ```

- Run specific tests for controllers:

  ```bash
  go test ./delivery/controller -v
  ```

- Run specific tests for middlewares:

  ```bash
  go test ./infrastructure -v
  ```

## Continuous Integration

A CI/CD pipeline is set up using GitHub Actions to automatically run tests on each push or pull request to the main branch. This ensures that the code is always in a healthy state.

Beware that the test for the repository layer is not part of the CI/CD pipeline since a test database is used to implement the testing for the repository layer

### This is the [API Documentation](https://documenter.getpostman.com/view/37363410/2sA3s3HB59) for this Simple Task Mangement System integrated with MongoDB with authentication/authorization features included
