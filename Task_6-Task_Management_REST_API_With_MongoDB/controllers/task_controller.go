package controllers

import (
	"errors"
	"net/http"

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

	return nil
}

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

func ValidatePassword(curr_password string, existing_password string) error {
	return bcrypt.CompareHashAndPassword([]byte(existing_password), []byte(curr_password))
}

func GenerateSignedToken(user_id string, user_email string, user_role string) (string, error) {
	// Generate JWT with claims of 'user_id', 'user_email', 'user_role'
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user_id,
		"email":     user_email,
		"user_role": user_role,
	})

	// sign the token with the secret key 'jwtSecret'
	return token.SignedString(data.JwtSecret) // [JWT token]string, [err]error
}

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

func HandleUserPromotion(c *gin.Context) {
	// get user role from the context
	user_role, err := GetUserRoleFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if user_role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		c.Abort()
		return
	}

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
	// get user role from the context
	user_role, err := GetUserRoleFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if user_role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		c.Abort()
		return
	}

	var updated_task models.AddedTask
	id := c.Param("id")

	if e := c.BindJSON(&updated_task); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = systemManager.UpdateTask(id, &updated_task)

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
	// get user role from the context
	user_role, err := GetUserRoleFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if user_role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		c.Abort()
		return
	}

	id := c.Param("id")

	err = systemManager.DeleteTask(id)

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
	// get user role from the context
	user_role, err := GetUserRoleFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if user_role != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		c.Abort()
		return
	}

	var new_task models.AddedTask

	if e := c.BindJSON(&new_task); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = systemManager.AddTask(&new_task)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task added successfully"})
}
