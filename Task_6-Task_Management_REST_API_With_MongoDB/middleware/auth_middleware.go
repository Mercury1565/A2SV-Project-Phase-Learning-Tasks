package middleware

import (
	"Task_6-Task_Management_REST_API_With_MongoDB/data"

	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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
		// the JWT token is valid (signatures verified and claims validated)
		token, err := jwt.Parse(authParts[1], TokenParser)
		if err != nil || !token.Valid {
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
