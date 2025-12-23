package handlers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oloomoses/opinions-hub/internal/dto"
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

	payload := c.PostForm("payload")

	if payload == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content cannot be blank"})
		return
	}

	var req dto.CreateOpinionRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload is required"})
		return
	}

	// start db transaction
	tx := h.repo.DB.Begin()

	if tx.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	opinion := models.Opinion{Content: req.Content}

	if err := tx.Create(&opinion).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// handle multipart files
	form, err := c.MultipartForm()

	if err == nil && form.File["images"] != nil {
		files := form.File["images"]

		var images []models.Image

		for _, file := range files {
			if !isImage(file) {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image type"})
				return
			}

			filename := generateFilename(file.Filename)
			path := filepath.Join("uploads/opinion_images", filename)

			if err := c.SaveUploadedFile(file, path); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the selected image"})
				return
			}

			images = append(images, models.Image{
				OpinionID: opinion.ID,
				URL:       path,
				MimeType:  file.Header.Get("Content-Type"),
				Size:      file.Size,
			})
		}

		if err := tx.Create(&images).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save images"})
			return
		}

		opinion.Images = images
	}

	tx.Commit()

	c.JSON(http.StatusCreated, opinion)
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

func isImage(file *multipart.FileHeader) bool {
	allowed := []string{"image/jpeg", "image/png", "image/webp", "image/jpg"}

	contentType := file.Header.Get("Content-Type")

	for _, t := range allowed {
		if t == contentType {
			return true
		}
	}
	return false
}

func generateFilename(original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}
