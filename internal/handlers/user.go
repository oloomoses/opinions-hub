package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oloomoses/opinions-hub/internal/dto"
	"github.com/oloomoses/opinions-hub/internal/models"
	"github.com/oloomoses/opinions-hub/internal/repository"
	"github.com/oloomoses/opinions-hub/internal/service/auth"
	"github.com/oloomoses/opinions-hub/internal/service/validator"
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

	if err := validator.ValidateCreateUser(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user := models.User{
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Username:     strings.TrimSpace(req.Username),
		PasswordHash: hashedPassword,
	}

	if err := h.repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userId, _ := c.Get("user_id")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"id":     userId,
		"handle": username,
	})
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

func (h *UserHandler) FollowUser(c *gin.Context) {
	currentUserID := c.GetUint("user_id")
	followingID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if currentUserID == uint(followingID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot self follow"})
		return
	}

	if err := h.repo.Follow(currentUserID, uint(followingID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "follow success"})

}

func (h *UserHandler) Unfollow(c *gin.Context) {
	currentUserID := c.GetUint("user_id")
	followingID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.repo.Unfollow(currentUserID, uint(followingID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unfollowing failed!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfollow success"})

}

func (h *UserHandler) GetFollowers(c *gin.Context) {
	currentUserID := c.GetUint("user_id")

	followers, err := h.repo.GetFollowers(currentUserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch followers, try again later!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"followers": followers})
}

func (h *UserHandler) GetFollowing(c *gin.Context) {
	currentUserID := c.GetUint("user_id")

	following, err := h.repo.GetFollowing(currentUserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch following, try again later!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": following})
}
