package repository

import (
	"Task_7-Clean_Architecture_Task_Management_REST_API/domain"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	database   mongo.Database
	collection string
}

func NewUserRepo(database mongo.Database, collection string) domain.UserRepository {
	return &userRepo{
		database:   database,
		collection: collection,
	}
}

// Create inserts a new user into the database.
// It takes a context and a user object as parameters.
// It returns an error if the insertion fails.
func (userRepo *userRepo) Create(c context.Context, user *domain.User) error {
	collection := userRepo.database.Collection(userRepo.collection)

	user.UserID = primitive.NewObjectID()
	_, err := collection.InsertOne(c, user)
	return err
}

// GetByEmail retrieves a user from the database by their email address.
// It takes a context and an email string as input parameters.
// It returns a pointer to a domain.User struct and an error.
// If the user is found, the function returns a pointer to the user struct and a nil error.
// If the user is not found, the function returns nil and a nil error.
// If an error occurs during the database operation, the function returns nil and the error.
func (userRepo *userRepo) GetByEmail(c context.Context, email string) (*domain.User, error) {
	collection := userRepo.database.Collection(userRepo.collection)

	var user domain.User
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user from the database based on the provided userID.
// It returns the user and any error encountered during the retrieval process.
func (userRepo *userRepo) GetByID(c context.Context, userID string) (*domain.User, error) {
	collection := userRepo.database.Collection(userRepo.collection)

	var user domain.User
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return &user, err
	}

	err = collection.FindOne(c, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, err
}

// UpdateUser updates the information of a user in the database.
// It takes a context and a user object as parameters.
// Returns an error if there was a problem updating the user.
func (userRepo *userRepo) UpdateUser(c context.Context, user *domain.User) error {
	collection := userRepo.database.Collection(userRepo.collection)

	filter := bson.M{"_id": user.UserID}
	update := bson.M{"$set": user}

	_, err := collection.UpdateOne(c, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// AreThereAnyUsers checks if there are any users in the database.
// It returns true if there are users, false otherwise.
// An error is returned if there was a problem counting the documents.
func (userRepo *userRepo) AreThereAnyUsers(c context.Context) (bool, error) {
	collection := userRepo.database.Collection(userRepo.collection)

	count, err := collection.CountDocuments(c, bson.M{})
	if err != nil {
		return false, nil
	}

	return count > 0, nil
}
