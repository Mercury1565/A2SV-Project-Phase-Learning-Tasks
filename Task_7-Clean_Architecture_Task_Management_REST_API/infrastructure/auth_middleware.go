package infrastructure

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.IndentedJSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		splitted := strings.Split(authHeader, " ")
		if len(splitted) != 2 || strings.ToLower(splitted[0]) != "bearer" {
			c.IndentedJSON(401, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := splitted[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			// return the secret key
			return []byte(secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
			c.Abort()
			return
		}

		// check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
			c.Abort()
			return
		}

		// set the claims to the context
		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

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
	user_role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("no user_role found in claims")
	}

	return user_role, nil
}
