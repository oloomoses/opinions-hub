package repository

import (
	"github.com/oloomoses/opinions-hub/internal/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(input *models.User) error {
	return r.DB.Create(&input).Error

}

func (r *UserRepo) UpdateUser(id uint64, updates map[string]interface{}) error {
	result := r.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepo) GetUser(id uint64) (models.User, error) {
	var user models.User

	err := r.DB.First(&user, id).Error

	return user, err
}

func (r *UserRepo) DeleteUser(id uint64) error {
	result := r.DB.Delete(&models.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}
