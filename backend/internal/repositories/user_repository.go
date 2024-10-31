package repositories

import (
	"backend/internal/models"

	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	GetUserByID(userID uint) (*models.User, error)
	DeleteUser(*models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(user *models.User) error {
	return r.db.Delete(user).Error
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}
