package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oloomoses/opinions-hub/internal/models"
	"github.com/oloomoses/opinions-hub/internal/repository"
)

type Opinion struct {
	repo *repository.Opinion
}

func NewOpinionHandler(repo *repository.Opinion) *Opinion {
	return &Opinion{repo: repo}
}

func (h *Opinion) CreateOpinion(c *gin.Context) {

	var input models.Opinion

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	if err := h.repo.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *Opinion) AllOpinions(c *gin.Context) {
	// var opnions []models.Opinion

	opnions, err := h.repo.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.JSON(http.StatusOK, opnions)
}

func (h *Opinion) UpdateOpinion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	var input struct {
		Content *string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	updates := make(map[string]interface{})

	updates["id"] = id

	if input.Content != nil {
		updates["content"] = *input.Content
	}

	if err := h.repo.Update(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	opinion, _ := h.repo.GetByID(uint(id))

	c.JSON(http.StatusOK, opinion)
}

func (h *Opinion) DeleteOpinion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
