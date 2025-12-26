package dto

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Username  string `json:"username" binding:"required,min=3,max=30"`
	Password  string `json:"password" binding:"required,min=4"`
}

type UserResponse struct {
	ID        uint64    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
