package repository

import (
	"errors"

	"github.com/oloomoses/opinions-hub/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func (r *UserRepo) VerifyUser(username string, password string) (models.User, error) {
	var user models.User

	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return user, errors.New("invald username or password!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return user, errors.New("invalid username or password")
	}

	return user, nil
}

func (r *UserRepo) Follow(followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("cannot self follow")
	}

	follower := &models.User{ID: int64(followerID)}

	result := r.DB.Model(follower).Association("Following").Append(&models.User{ID: int64(followingID)})

	return result
}

func (r *UserRepo) Unfollow(followerID, followingID uint) error {
	follower := &models.User{ID: int64(followerID)}

	return r.DB.Model(follower).Association("Following").Delete(&models.User{ID: int64(followingID)})
}

func (r *UserRepo) GetFollowers(userID uint) ([]models.User, error) {
	var user models.User

	err := r.DB.Preload("Followers").First(&user, userID).Error

	return user.Followers, err
}

func (r *UserRepo) GetFollowing(userID uint) ([]models.User, error) {
	var user models.User

	err := r.DB.Preload("Following").First(&user, userID).Error

	return user.Following, err
}

func (r *UserRepo) IsFollowing(followeID, followingID uint) (bool, error) {
	var count int64

	err := r.DB.Table("user_follows").
		Where("follower_id = ? AND following_id = ?", followeID, followingID).
		Count(&count).Error

	return count > 0, err
}
