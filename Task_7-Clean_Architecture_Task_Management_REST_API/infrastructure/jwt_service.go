package infrastructure

import (
	"Task_7-Clean_Architecture_Task_Management_REST_API/domain"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
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
func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// ExtractInfoFromToken extracts information from a JWT token.
// It takes a requestToken string and a secret string as input parameters.
// It returns the extracted user ID and role as strings, along with any error encountered.
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
