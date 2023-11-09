package utils

import (
	"errors"

	bcryptL "golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4
	MaxCost     int = 31
	DefaultCost int = 10
)

// Compare  Bcrypt.
func ComparePassword(hashedPassword string, password string) (bool, error) {
	err := bcryptL.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcryptL.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Hash  Bcrypt.
func HashPassword(password string, options ...func(*hashOptions)) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password must be less than 8 characters")
	}
	opts := hashOptions{cost: DefaultCost}
	for _, opt := range options {
		opt(&opts)
	}
	if opts.cost <= MinCost || opts.cost >= MaxCost {
		return "", errors.New("cost must be between 4 and 31")
	}
	hashPassword, err := bcryptL.GenerateFromPassword([]byte(password), opts.cost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

type hashOptions struct {
	cost int
}
