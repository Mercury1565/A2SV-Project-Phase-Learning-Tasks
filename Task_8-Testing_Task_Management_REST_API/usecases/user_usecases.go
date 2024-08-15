package usecases

import (
	"Task_8-Testing_Task_Management_REST_API/domain"
	"Task_8-Testing_Task_Management_REST_API/infrastructure"
	"context"
	"time"
)

type userUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (userUC *userUsecase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, userUC.contextTimeout)
	defer cancel()
	return userUC.userRepository.Create(ctx, user)
}

func (userUC *userUsecase) GetByEmail(c context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, userUC.contextTimeout)
	defer cancel()
	return userUC.userRepository.GetByEmail(ctx, email)
}

func (userUC *userUsecase) GetByID(c context.Context, userID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, userUC.contextTimeout)
	defer cancel()
	return userUC.userRepository.GetByID(ctx, userID)
}

func (userUC *userUsecase) UpdateUser(c context.Context, updated_user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, userUC.contextTimeout)
	defer cancel()
	return userUC.userRepository.UpdateUser(ctx, updated_user)
}

func (userUC *userUsecase) AreThereAnyUsers(c context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(c, userUC.contextTimeout)
	defer cancel()
	return userUC.userRepository.AreThereAnyUsers(ctx)
}

func (loginUsecase *userUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (string, error) {
	return infrastructure.CreateAccessToken(user, secret, expiry)
}
