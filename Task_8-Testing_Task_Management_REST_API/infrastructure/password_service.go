package infrastructure

import "golang.org/x/crypto/bcrypt"

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
