# Task Mangment System with GO and MongoDB

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
   cd Task_5-Task_Management_REST_API_With_MongoDB
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

## Running the Project

1. Run the project:

   ```bash
   go run main.go
   ```

2. The application will be running at `http://localhost:8080`. You can change the application port number in the `main.go` file

## API Endpoints

- GET Requests

  - http://localhost:8080/tasks : Get tasks
  - http://localhost:8080/tasks/ID : Get task with taskId ID

- PUT Request

  - http://localhost:8080/tasks/ID: Update the fields of task with taskId ID

- DELETE Request

  - http://localhost:8080/tasks/ID: Delete the task with taskId ID

- POST Request
  - http://localhost:8080/tasks: Add new task

### This is the [API Documentation](https://documenter.getpostman.com/view/37363410/2sA3s1nrqe) for this Simple Task Mangement System integrated with MongoDB
