package validator

import (
	"errors"
	"strings"

	"github.com/oloomoses/opinions-hub/internal/dto"
)

func ValidateCreateUser(userReq dto.CreateUserRequest) error {
	var err error

	switch {
	case strings.TrimSpace(userReq.FirstName) == "":
		err = errors.New("First Name cannot be empty")
	case strings.TrimSpace(userReq.LastName) == "":
		err = errors.New("Last Name cannot be empty")
	case strings.TrimSpace(userReq.Username) == "":
		err = errors.New("username cannot be empty")
	default:
		err = nil
	}

	return err
}
