package infrastructure

import (
	"Task_8-Testing_Task_Management_REST_API/domain"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CreateAccessToken generates a JWT access token for the given user with the specified secret and expiry time.
// It takes a pointer to a User object, the secret key used for signing the token, and the expiry time in seconds.
// The function returns the generated access token as a string and any error encountered during the process.
func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()

	claims := &domain.JWTCustomClaims{
		Name: user.Name,
		ID:   user.UserID.Hex(),
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return signedToken, err
}

// IsAuthorized checks if the provided request token is authorized using the given secret.
// It returns a boolean indicating whether the token is authorized and an error if any.
func IsAuthorized(requestToken string, secret string) (*jwt.Token, error) {
	authorizedToken, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	return authorizedToken, err
}

// ExtractInfoFromToken extracts information from a JWT token.
// It takes a requestToken string and a secret string as input parameters.
// It returns the extracted user ID and role as strings, along with any error encountered.
// The case where the token is expired is implicitly handled in jwt.Parse. It returns a nonEmpty token with valid field equals to false
func ExtractInfoFromToken(requestToken string, secret string) (string, string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}

	return claims["id"].(string), claims["role"].(string), nil
}

// GetUserRoleFromContext retrieves the user role from the provided Gin context.
// It expects the context to contain a "claims" key, which should be a jwt.MapClaims object.
// If the "claims" key is not found or if the claims are not valid, an error is returned.
// If the user_role is not found in the claims or if it is not a string, an error is returned.
// Otherwise, the user_role is returned along with a nil error.
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
