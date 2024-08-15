package infrastructure

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware is a middleware function that performs JWT authentication.
// It checks the Authorization header for a valid JWT token and sets the claims to the context.
// If the token is invalid or missing, it returns an error response.
// The secret parameter is used to validate the token's signature.
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		splitted := strings.Split(authHeader, " ")
		if len(splitted) != 2 || strings.ToLower(splitted[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		// check if token is authorized
		tokenString := splitted[1]
		authorizedToken, err := IsAuthorized(tokenString, secret)
		// check if the token is valid
		if authorizedToken != nil && !authorizedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
			c.Abort()
			return
		}

		// set the claims to the context
		claims := authorizedToken.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		c.Next()
	}
}
