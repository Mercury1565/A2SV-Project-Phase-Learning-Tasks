package repository

import (
	"Task_8-Testing_Task_Management_REST_API/domain"
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepo struct {
	database   mongo.Database
	collection string
}

func NewTaskRepo(database mongo.Database, collection string) domain.TaskRepository {
	return &taskRepo{
		database:   database,
		collection: collection,
	}
}

// Create inserts a new task into the database.
// It takes a context and a task object as parameters.
// It returns an error if the insertion fails.
func (taskRepo *taskRepo) Create(c context.Context, task *domain.Task) error {
	collection := taskRepo.database.Collection(taskRepo.collection)

	task.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(c, task)
	return err
}

// GetTasks retrieves all tasks from the database.
// It returns a slice of domain.Task and an error, if any.
func (taskRepo *taskRepo) GetTasks(c context.Context) ([]domain.Task, error) {
	collection := taskRepo.database.Collection(taskRepo.collection)

	var tasks []domain.Task
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		return tasks, err
	}

	err = cursor.All(c, &tasks)
	if tasks == nil {
		return []domain.Task{}, err
	}

	return tasks, err
}

// GetTaskByID retrieves a task from the database based on the given task ID.
// It returns the retrieved task and an error, if any.
func (taskRepo *taskRepo) GetTaskByID(c context.Context, taskID string) (domain.Task, error) {
	collection := taskRepo.database.Collection(taskRepo.collection)

	var task domain.Task
	obj_ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return task, err
	}

	err = collection.FindOne(c, bson.M{"_id": obj_ID}).Decode(&task)
	if err != nil {
		return task, err
	}

	return task, err
}

// UpdateTask updates a task with the specified taskID in the repository.
// It takes a context, taskID string, and updated_task *domain.Task as parameters.
// The function returns an error if any occurred during the update process.
func (taskRepo *taskRepo) UpdateTask(c context.Context, taskID string, updated_task *domain.Task) error {
	collection := taskRepo.database.Collection(taskRepo.collection)

	obj_ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	updated_fields := make(bson.M)

	// populate the update parameter by checking validity of the new task
	if updated_task.Title != "" {
		updated_fields["title"] = updated_task.Title
	}
	if updated_task.Description != "" {
		updated_fields["description"] = updated_task.Description
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

	// update the task with id 'taskID'
	updateResult, err := collection.UpdateOne(c, filter, update)

	if err != nil {
		log.Fatal(err)
	}

	if updateResult.MatchedCount == 0 {
		// task with id 'taskID' not found
		return fmt.Errorf("task with id '%v' not found", taskID)
	}

	return nil
}

// DeleteTask deletes a task from the repository based on the given task ID.
// It takes a context `c` and a task ID `taskID` as parameters.
// It returns an error if the task ID is invalid or if the task with the given ID is not found.
func (taskRepo *taskRepo) DeleteTask(c context.Context, taskID string) error {
	collection := taskRepo.database.Collection(taskRepo.collection)

	obj_ID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return errors.New("invalid id entered")
	}

	// define filter
	filter := bson.M{
		"_id": obj_ID,
	}

	// delete task with id 'taskID'
	deleteResult, err := collection.DeleteOne(c, filter)

	if err != nil {
		log.Fatal(err)
	}

	if deleteResult.DeletedCount == 0 {
		// task with id 'taskID' not found
		return fmt.Errorf("task with id '%v' not found", taskID)
	}

	return nil
}
