package infrastructure

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

func ValidatePassword(curr_password string, existing_password string) error {
	return bcrypt.CompareHashAndPassword([]byte(existing_password), []byte(curr_password))
}
