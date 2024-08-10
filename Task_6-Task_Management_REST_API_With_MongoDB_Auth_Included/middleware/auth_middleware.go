package middleware

import (
	"Task_6-Task_Management_REST_API_With_MongoDB/data"
	"errors"
	"net/http"
	"time"

	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// TokenParser is a function that parses a JWT token and returns the secret key used for signing the token.
// It checks if the signing method of the token is the expected one and returns an error if it's not.
// The function takes a pointer to a jwt.Token as input and returns the secret key and an error.
func TokenParser(token *jwt.Token) (interface{}, error) {
	// check if the signing method of the token is one we are looking for
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return data.JwtSecret, nil
}

// this is where the authorization action is performed
func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.IndentedJSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// convert the auth header to a convinient string type slice 'authParts'
		// expected Auth-Header -> Bearer <token>
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.IndentedJSON(401, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		// check if...
		// no error occured during parsing
		// the JWT hasn't expired
		// the JWT token is valid (signatures verified and claims validated)
		token, err := jwt.Parse(authParts[1], TokenParser)
		if err != nil || !token.Valid {
			c.IndentedJSON(401, gin.H{"error": "invalid JWT"})
			c.Abort()
			return
		}

		// check if the token is expired
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			exp, ok := claims["expiration"].(float64)
			if !ok || time.Now().Unix() > int64(exp) {
				// JWT is expired
				c.IndentedJSON(401, gin.H{"error": "JWT is expired"})
				c.Abort()
				return
			}
		} else {
			// JWT is invalid for other reasons
			c.IndentedJSON(401, gin.H{"error": "invalid JWT"})
			c.Abort()
			return
		}

		// store the claims in the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// store the parsed JWT in the gin context
			c.Set("claims", claims)
		}

		c.Next()
	}
}

func AuthenticateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
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
