# Task Mangment System with GO

## Prerequisites

- [Go](https://golang.org/doc/install)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/Mercury1565/A2SV-Project-Phase-Learning-Tasks
   ```

2. Navigate to the project directory:

   ```bash
   cd Task_4-Task_Management_REST_API
   ```

3. Download the required Go dependencies:

   ```bash
   go mod download
   ```

4. Tidy up the dependencies:
   ```bash
   go mod tidy
   ```

## Running the Project

1. Run the project:

   ```bash
   go run main.go
   ```

2. The application will be running at `http://localhost:8080`.

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

### This is the [API Documentation](https://documenter.getpostman.com/view/37363410/2sA3rzJCHr) for this Simple Task Mangement System.
