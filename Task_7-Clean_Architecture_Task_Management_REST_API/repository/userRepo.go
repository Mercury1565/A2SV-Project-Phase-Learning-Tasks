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

func (userRepo *userRepo) Create(c context.Context, user *domain.User) error {
	collection := userRepo.database.Collection(userRepo.collection)

	user.UserID = primitive.NewObjectID()
	_, err := collection.InsertOne(c, user)
	return err
}

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

func (userRepo *userRepo) AreThereAnyUsers(c context.Context) (bool, error) {
	collection := userRepo.database.Collection(userRepo.collection)

	count, err := collection.CountDocuments(c, bson.M{})
	if err != nil {
		return false, nil
	}

	return count > 0, nil
}
