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

// ValidateUserInfo validates the user information before performing any operations.
// It checks if the password is at least 6 characters long, if the user role is either 'USER' or 'ADMIN',
// and if the name field is not empty. It also checks if there are any existing users in the system.
// If the user role is 'ADMIN' and there are existing users, it returns an error indicating that
// an admin can only be registered if no users exist.
// If all validations pass, it returns nil.
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

	if user.Role == "ADMIN" && usersExist {
		fmt.Println(usersExist)
		return errors.New("admin can only be registered if no users exist")
	}

	return nil
}

// HandelUserRegister handles the registration of a new user.
// It receives the user registration information from the request body,
// validates the information, checks if the user already exists,
// hashes the password, and adds the user to the database.
// If successful, it returns a JSON response with a success message.
// If there are any errors, it returns a JSON response with an error message.
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

// HandelUserLogin handles the user login functionality.
// It receives a request context and expects the user information to be provided in the request body as JSON.
// It checks if the user exists and if the provided password is correct.
// If the user exists and the password is correct, it generates a signed JWT token and returns it in the response.
// The token can be used for authentication in subsequent requests.
// If there are any errors during the process, appropriate error responses are returned.
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
		context.JSON(401, gin.H{"error": "incorrect password"})
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

// HandleUserPromotion handles the promotion of a user to the 'ADMIN' role.
// It takes a gin.Context object as a parameter and retrieves the user ID from the request parameters.
// It then checks if the user exists and if they are already an admin.
// If the user is not found, it returns a JSON response with an error message.
// If the user is already an admin, it returns a JSON response indicating that the user is already an admin.
// If the user is not an admin, it promotes the user to the 'ADMIN' role and saves the changes in the database.
// If there is an error during the update, it returns a JSON response with an error message.
// Finally, it returns a JSON response indicating that the user has been promoted to admin status.
func (controller *UserController) HandleUserPromotion(c *gin.Context) {
	id := c.Param("id")

	existingUser, err := controller.UserUsecase.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if existingUser.Role == "ADMIN" {
		c.JSON(http.StatusOK, gin.H{"message": "user is already an admin"})
		return
	}

	// promote user to 'ADMIN'
	existingUser.Role = "ADMIN"

	// save the changes in the database
	err = controller.UserUsecase.UpdateUser(c, existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin status"})
}
