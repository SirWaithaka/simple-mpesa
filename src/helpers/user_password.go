package helpers

import (
	"simple-mpesa/src/errors"

	"golang.org/x/crypto/bcrypt"
)

// ComparePasswordToHash verify hashed password and plain text password if they match.
func ComparePasswordToHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.Error{Code: errors.EINVALID, Message: errors.InvalidCredentials}
		}
		return errors.PasswordHashError{Err: err}
	}

	return nil
}

// HashPassword takes a plain text password and hash it.
func HashPassword(password string) (string, error) {
	// generate hashed password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.PasswordHashError{Err: err}
	}

	return string(hash), nil
}
