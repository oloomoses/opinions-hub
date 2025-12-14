package repository

import (
	"github.com/oloomoses/opinions-hub/internal/models"
	"gorm.io/gorm"
)

type Opinion struct {
	DB *gorm.DB
}

func NewOpinionRepo(db *gorm.DB) *Opinion {
	return &Opinion{DB: db}
}

func (op *Opinion) Create(content string) error {
	return op.DB.Create(content).Error
}

func (op *Opinion) GetAll() ([]models.Opinion, error) {
	var opinions []models.Opinion

	err := op.DB.Find(&opinions).Error

	return opinions, err
}

func (op *Opinion) GetByID(id uint) (models.Opinion, error) {
	var opinion models.Opinion

	err := op.DB.First(&opinion, id).Error

	return opinion, err
}

func (op *Opinion) Update(id uint, updates map[string]interface{}) error {
	return op.DB.Model(&models.Opinion{}).Where("id=?", id).Updates(updates).Error
}

func (op *Opinion) Delete(id uint) error {
	return op.DB.Delete(&models.Opinion{}, id).Error
}
