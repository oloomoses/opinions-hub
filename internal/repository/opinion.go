package repository

import (
	"github.com/oloomoses/opinions-hub/internal/models"
	"gorm.io/gorm"
)

type OpinionRepo interface {
	Create(*models.Opinion) error
	GetAll() ([]models.Opinion, error)
	GetByID(id uint) (models.Opinion, error)
	Update(id uint, update map[string]interface{}) error
	Delete(id uint) error
}

type Opinion struct {
	DB *gorm.DB
}

func NewOpinionRepo(db *gorm.DB) *Opinion {
	return &Opinion{DB: db}
}

func (op *Opinion) Create(opinion *models.Opinion) error {
	if err := opinion.ValidateContentCannotBeEmpty(); err != nil {
		return err
	}

	return op.DB.Create(opinion).Error
}

func (op *Opinion) GetAll() ([]models.Opinion, error) {
	var opinions []models.Opinion

	err := op.DB.Order("id ASC").Find(&opinions).Error

	return opinions, err
}

func (op *Opinion) GetByID(id uint) (models.Opinion, error) {
	var opinion models.Opinion

	err := op.DB.First(&opinion, id).Error

	return opinion, err
}

func (op *Opinion) Update(id uint, updates map[string]interface{}) error {
	result := op.DB.Model(&models.Opinion{}).Where("id=?", id).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (op *Opinion) Delete(id uint) error {
	result := op.DB.Delete(&models.Opinion{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
