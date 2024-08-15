package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionUser = "users"

type User struct {
	UserID   primitive.ObjectID `json:"-" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	GetByEmail(c context.Context, email string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
	UpdateUser(c context.Context, user *User) error
	AreThereAnyUsers(c context.Context) (bool, error)
}

type UserUsecase interface {
	Create(c context.Context, user *User) error
	GetByEmail(c context.Context, email string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
	UpdateUser(c context.Context, user *User) error
	AreThereAnyUsers(c context.Context) (bool, error)
	CreateAccessToken(user *User, secret string, expiry int) (string, error)
}
