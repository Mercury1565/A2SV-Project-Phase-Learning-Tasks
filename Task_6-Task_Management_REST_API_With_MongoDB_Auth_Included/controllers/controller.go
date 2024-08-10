package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"Task_6-Task_Management_REST_API_With_MongoDB/data"
	"Task_6-Task_Management_REST_API_With_MongoDB/models"
)

// declare systemManager instance
var systemManager *data.SystemManagement

func InitializeSystemMangement() {
	// initialize a systemManager instance
	systemManager = data.NewSystemManager()
}

// ValidateUserInfo validates the user information provided.
// It checks if the password is at least 6 characters long,
// if the user role is either 'USER' or 'ADMIN',
// and if the name field is not empty.
// It also checks if there are any existing users in the system.
// If the user role is 'ADMIN' and there are existing users,
// it returns an error indicating that an admin can only be registered if no users exist.
// Otherwise, it returns nil if all validations pass.
func ValidateUserInfo(user *models.NewUser) error {
	if len(user.Password) < 6 {
		return errors.New("password must be atleast 6 characters long")
	}
	if user.Role != "ADMIN" && user.Role != "USER" {
		return errors.New("invalid user role, user role is either 'USER' or 'ADMIN'")
	}
	if len(user.Name) == 0 {
		return errors.New("empty name field not allowed")
	}

	usersExist, err := systemManager.AreThereAnyUsers()
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

// HashPassword hashes the given password using bcrypt algorithm.
// It returns the hashed password as a byte slice and any error encountered during the hashing process.
func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

// ValidatePassword compares the current password with the existing password hash.
// It uses bcrypt.CompareHashAndPassword to perform the comparison.
// If the passwords match, it returns nil. Otherwise, it returns an error.
func ValidatePassword(curr_password string, existing_password string) error {
	return bcrypt.CompareHashAndPassword([]byte(existing_password), []byte(curr_password))
}

// GenerateSignedToken generates a signed JWT token with the provided user information.
// It takes the user ID, user email, and user role as input parameters and returns the generated token as a string.
// The token is signed using the secret key 'jwtSecret' and includes claims for 'user_id', 'user_email', 'user_role', and 'expiration'.
// The expiration time of the token is set to 24 hours from the current time.
// If an error occurs during token generation, it is returned as the second return value.
func GenerateSignedToken(user_id string, user_email string, user_role string) (string, error) {
	// set the expiration time of the token to 24Hrs
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Generate JWT with claims of 'user_id', 'user_email', 'user_role', 'exp_time'
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user_id,
		"email":      user_email,
		"user_role":  user_role,
		"expiration": expirationTime,
	})

	// sign the token with the secret key 'jwtSecret'
	return token.SignedString(data.JwtSecret) // [JWT token]string, [err]error
}

// GetUserRoleFromContext retrieves the user role from the provided Gin context.
// It expects the context to contain a "claims" key, which should be a jwt.MapClaims object.
// The function returns the user role as a string and an error if any.
func GetUserRoleFromContext(context *gin.Context) (string, error) {
	// retrieve claims from the context
	claimsValue, exists := context.Get("claims")

	if !exists {
		return "", errors.New("no claims found")
	}

	// retrieve jwt.MapClaims from the claimsValue
	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
		return "", errors.New("claims are not valid")
	}

	// retrieve the user_role from the claims
	user_role, ok := claims["user_role"].(string)
	if !ok {
		return "", errors.New("no user_role found in claims")
	}

	return user_role, nil
}

// HandelUserRegister handles the registration of a new user.
// It receives the user information from the request body and performs the following steps:
// 1. Validates the user information.
// 2. Checks if the user already exists in the database.
// 3. Hashes the user's password.
// 4. Adds the new user to the database.
func HandelUserRegister(context *gin.Context) {
	var curr_user models.NewUser

	// get inputed info from body
	err := context.ShouldBindJSON(&curr_user)
	if err != nil {
		context.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// check validity of the info entered entered for the new user
	err = ValidateUserInfo(&curr_user)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	// check if user already exists
	existingUser, err := systemManager.GetExistingUserByEmail(curr_user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingUser != nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": "user already exists"})
		return
	}

	// hash inputed password
	hashedPassword, err := HashPassword(curr_user.Password)
	if err != nil {
		context.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	curr_user.Password = string(hashedPassword)

	// add user to database
	err = systemManager.AddNewUser(&curr_user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "user registered successfully"})
}

// HandelUserLogin handles the user login process.
// It receives the user login data from the request body and validates it.
// If the user exists and the password is correct, it generates a signed JWT token and returns it in the response.
// The token contains the user's ID, email, and role as claims.
// If any error occurs during the process, it returns an appropriate error response.
func HandelUserLogin(context *gin.Context) {
	var curr_user *models.LoginData

	// get inputed info from body
	err := context.ShouldBindJSON(&curr_user)
	if err != nil {
		context.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// fetch user
	existingUser, err := systemManager.GetExistingUserByEmail(curr_user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingUser == nil {
		context.JSON(http.StatusNotAcceptable, gin.H{"error": "user doesn't exists"})
		return
	}

	// check if user has inputed the correct password
	err = ValidatePassword(curr_user.Password, existingUser.Password)
	if err != nil {
		context.JSON(401, gin.H{"error": "Incorrect password"})
		return
	}

	// generate signed JWT with 'user_id', 'user_email' and 'user_role' claims
	signed_jwt_token, err := GenerateSignedToken(existingUser.UserID.String(), existingUser.Email, existingUser.Role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	context.JSON(200, gin.H{"message": "user logged in successfully", "token": signed_jwt_token})
}

// HandleUserPromotion handles the promotion of a user to admin status.
// It takes a gin.Context object as a parameter and retrieves the user ID from the request parameters.
// It then converts the ID to an object ID and retrieves the existing user from the system manager.
// If the user is not found, it returns a 404 error.
// If the user is already an admin, it returns a success message indicating user is already an ADMIN
// Otherwise, it promotes the user to admin status, saves the changes in the database, and returns a success message.
func HandleUserPromotion(c *gin.Context) {
	id := c.Param("id")

	objId, err := data.ConvertToObjectID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	existingUser, err := systemManager.GetExistingUserByID(objId)
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

	// Save the changes in the database
	err = systemManager.UpdateUser(existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin status"})

}

// GetAllTasks returns all tasks.
func GetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, systemManager.GetAllTasks())
}

// GetTask retrieves a task by its ID from the task collection and returns it as JSON.
// If the task is not found, it returns a JSON response with a "task not found" message.
func GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := systemManager.GetTask(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// UpdateTask updates a task with the given ID.
// It expects a JSON request body containing the updated task details.
// If the request body is invalid, it returns a 400 Bad Request error with a "invalid request body" error message.
// If the task is not found, it returns a 404 Not Found error with a "task not found" message.
// Otherwise, it returns a 200 OK response with a success message.
func UpdateTask(c *gin.Context) {
	var updated_task models.AddedTask
	id := c.Param("id")

	if e := c.BindJSON(&updated_task); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := systemManager.UpdateTask(id, &updated_task)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// DeleteTask deletes a task with the given ID.
// If the task is not found, it returns a 404 Not Found error with a "task not found" message.
// Otherwise, it returns a 200 OK response with a success message.
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := systemManager.DeleteTask(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

// CreateTask creates a new task based on the JSON data provided in the request body.
// It binds the JSON data to the `new_task` variable and adds it to the task collection.
// If the request body is invalid, it returns a JSON response with an error message.
// If the task is added successfully, it returns a JSON response with a success message.
func CreateTask(c *gin.Context) {
	var new_task models.AddedTask

	if e := c.BindJSON(&new_task); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := systemManager.AddTask(&new_task)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task added successfully"})
}
