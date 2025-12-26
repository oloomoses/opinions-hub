package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oloomoses/opinions-hub/internal/dto"
	"github.com/oloomoses/opinions-hub/internal/models"
	"github.com/oloomoses/opinions-hub/internal/repository"
	"github.com/oloomoses/opinions-hub/internal/service/auth"
)

type UserHandler struct {
	repo *repository.UserRepo
}

func NewUserHandler(repo *repository.UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := dto.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Username:     req.Username,
		PasswordHash: hashedPassword,
	}

	if err := h.repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var input dto.LoginUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := strings.TrimSpace(input.Username)
	password := input.Password

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password cannot be empty"})
		return
	}

	user, err := h.repo.VerifyUser(username, password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, _ := auth.GenerateAccessToken(uint(user.ID), user.Username)

	c.JSON(http.StatusOK, gin.H{
		"User":        user,
		"AccessToken": accessToken,
	})

}
