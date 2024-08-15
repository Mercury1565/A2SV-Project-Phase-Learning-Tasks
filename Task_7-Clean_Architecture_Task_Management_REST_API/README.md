# Task Mangment System with GO | MongoDB | Authentication & Authorizaiton Features Included
## This repository shows how clean architecture is implemented in software system using a simple Task Mangement App

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
   cd Task_7-Clean_Architecture_Task_Management_REST_API
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
├── delivery/
│   ├── controllers/
│   │   ├── user_controller.go
│   │   └── tasks_controller.go
│   ├── routers/
│   │   ├── adminRoutes.go
│   │   ├── protectedRoutes.go
│   │   ├── publicRoutes.go
│   │   └── route.go
│   ├── .env
│   └── main.go
├── domain/
│   ├── jwtCustom.go
│   ├── task.go
│   └── user.go
├── infrastructure/
│   ├── auth_middleware.go
│   ├── jwt_service.go
│   └── password_service.go
├── repository/
│   ├── taskRepo.go
│   └── userRepo.go
├── usecases/
│   ├── user_usecases.go
│   └── task_usecases.go
├── README.mod
├── go.mod
└── go.sum
```

## Running the Project

1. Navigate to the delivery directory:

   ```bash
   cd delivery
   ```

2. Run the project:

   ```bash
   go run main.go
   ```

3. The application will be running at `http://localhost:8080`. You can change the application port number as well us other environment variables in the `.env` file located in teh `/delivery` directory

## API Endpoints

### APIs Related to Authentication

- POST Requests

  - http://localhost:8080/register: register new user
  - http://localhost:8080/login : authenticate and signin users
  - http://localhost:8080/promote/userID : promote user role to 'ADMIN', only allowed for users with 'ADMIN'

### APIs Related to task managment

- GET Requests

  - http://localhost:8080/tasks : Get tasks
  - http://localhost:8080/tasks/taskID : Get task with taskId ID

- PUT Request

  - http://localhost:8080/tasks/taskID: Update the fields of task with taskId ID, only allowed for users with 'ADMIN'

- DELETE Request

  - http://localhost:8080/tasks/taskID: Delete the task with taskId ID, only allowed for users with 'ADMIN'

- POST Request

  - http://localhost:8080/tasks: Add new task, only allowed for users with 'ADMIN'

### This is the [API Documentation](https://documenter.getpostman.com/view/37363410/2sA3s3HB59) for this Simple Task Mangement System integrated with MongoDB with authentication/authorization features included
