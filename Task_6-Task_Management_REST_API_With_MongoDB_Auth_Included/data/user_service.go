package data

import (
	"Task_6-Task_Management_REST_API_With_MongoDB/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Global variable to store the JWT secret
var JwtSecret = []byte("your_jwt_secret")

// define mongoDB task_collection for tasks
var user_collection *mongo.Collection

// AddNewUser adds a new user to the system.
// It takes a pointer to a `NewUser` struct as a parameter and returns an error.
// The function inserts the new user into the user collection in the database.
// If an error occurs during the insertion, it returns a `server error` error.
func (userManager *SystemManagement) AddNewUser(new_user *models.NewUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userManager.user_collection.InsertOne(ctx, new_user)

	if err != nil {
		return errors.New("server error")
	}

	return nil
}

// AreThereAnyUsers checks if there are any users in the system.
// It returns true if there are users, false otherwise.
// An error is returned if there was a problem checking for users.
func (userManager *SystemManagement) AreThereAnyUsers() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userManager.user_collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, nil
	}

	return count > 0, nil
}

// GetExistingUserByEmail retrieves an existing user from the database based on the provided email.
// If the user is found, it returns a pointer to the User object and nil error.
// If the user is not found, it returns nil and nil error.
// If any error occurs during the retrieval process, it returns nil and the error.
func (userManager *SystemManagement) GetExistingUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := userManager.user_collection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &existingUser, nil
}

// GetExistingUserByID retrieves an existing user from the database based on the provided ID.
// If the user is found, it returns a pointer to the User object and nil error.
// If the user is not found, it returns nil and nil error.
// If any other error occurs during the retrieval process, it returns nil and the error.
func (userManager *SystemManagement) GetExistingUserByID(id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := userManager.user_collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &existingUser, nil
}

// GetUser retrieves a user from the user collection based on the provided email.
// If the user is found, it returns an error indicating that the user already exists.
// If the user is not found, it returns nil.
func (userManager *SystemManagement) GetUser(new_user *models.NewUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userManager.user_collection.FindOne(ctx, bson.M{"email": new_user.Email}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Fatal(err)
	}
	return fmt.Errorf("user with email %s already exists", new_user.Email)
}

// UpdateUser updates the information of a user in the system.
// It takes a pointer to a User struct as a parameter and returns an error if any.
func (userManager *SystemManagement) UpdateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.UserID}
	update := bson.M{"$set": user}

	_, err := userManager.user_collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
