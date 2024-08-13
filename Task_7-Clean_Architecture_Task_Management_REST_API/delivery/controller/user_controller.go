package controller

import (
	"Task_7-Clean_Architecture_Task_Management_REST_API/bootstrap"
	"Task_7-Clean_Architecture_Task_Management_REST_API/domain"
	"Task_7-Clean_Architecture_Task_Management_REST_API/infrastructure"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase domain.UserUsecase
	Env         *bootstrap.Env
}

func (controller *UserController) ValidateUserInfo(c context.Context, user *domain.User) error {
	if len(user.Password) < 6 {
		return errors.New("password must be atleast 6 characters long")
	}
	if user.Role != "ADMIN" && user.Role != "USER" {
		return errors.New("invalid user role, user role is either 'USER' or 'ADMIN'")
	}
	if len(user.Name) == 0 {
		return errors.New("empty name field not allowed")
	}

	usersExist, err := controller.UserUsecase.AreThereAnyUsers(c)
	if err != nil {
		return err
	}

	fmt.Println(user.Role, usersExist)

	if user.Role == "ADMIN" && usersExist {
		fmt.Println(usersExist)
		return errors.New("admin can only be registered if no users exist")
	}

	return nil
}

func (controller *UserController) HandelUserRegister(context *gin.Context) {
	var curr_user *domain.User

	// get inputed info from body
	err := context.ShouldBindJSON(&curr_user)
	if err != nil {
		context.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// check validity of the info entered entered for the new user
	err = controller.ValidateUserInfo(context, curr_user)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	// check if user already exists
	existingUser, err := controller.UserUsecase.GetByEmail(context, curr_user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingUser != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": "user already exists"})
		return
	}

	// hash inputed password
	hashedPassword, err := infrastructure.HashPassword(curr_user.Password)
	if err != nil {
		context.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	curr_user.Password = string(hashedPassword)

	// add user to database
	err = controller.UserUsecase.Create(context, curr_user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "user registered successfully"})
}

func (controller *UserController) HandelUserLogin(context *gin.Context) {
	var curr_user *domain.User

	// get inputed info from body
	err := context.ShouldBindJSON(&curr_user)
	if err != nil {
		context.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// check if user already exists
	existingUser, err := controller.UserUsecase.GetByEmail(context, curr_user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingUser == nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": "user doesnt't exists"})
		return
	}

	// check if user has inputed the correct password
	err = infrastructure.ValidatePassword(curr_user.Password, existingUser.Password)
	if err != nil {
		context.JSON(401, gin.H{"error": "Incorrect password"})
		return
	}

	accessTokenExp := controller.Env.AccessTokenExpiryHour
	accessTokenSecret := controller.Env.AccessTokenSecret

	// generate signed JWT with 'user_id', 'user_email' and 'user_role' claims
	signed_jwt_token, err := controller.UserUsecase.CreateAccessToken(existingUser, accessTokenSecret, accessTokenExp)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	context.JSON(200, gin.H{"message": "user logged in successfully", "token": signed_jwt_token})
}

func (controller *UserController) HandleUserPromotion(c *gin.Context) {
	id := c.Param("id")

	existingUser, err := controller.UserUsecase.GetByID(c, id)
	if err != nil {
		fmt.Println(1)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingUser == nil {
		fmt.Println(2)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if existingUser.Role == "ADMIN" {
		fmt.Println(3)
		c.JSON(http.StatusOK, gin.H{"message": "user is already an admin"})
		return
	}

	// promote user to 'ADMIN'
	existingUser.Role = "ADMIN"

	// Save the changes in the database
	err = controller.UserUsecase.UpdateUser(c, existingUser)
	if err != nil {
		fmt.Println(4)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin status"})
}
