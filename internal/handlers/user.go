package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oloomoses/opinions-hub/internal/dto"
	"github.com/oloomoses/opinions-hub/internal/models"
	"github.com/oloomoses/opinions-hub/internal/repository"
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

func (h *UserHandler) Login(c *gin.Context) {
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

	if err := h.repo.VerifyUser(username, password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)

}
