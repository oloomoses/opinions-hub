package handlers

import (
	"net/http"

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

	// h.repo.Update(id, updates)
}
