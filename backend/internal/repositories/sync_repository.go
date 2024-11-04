package repositories

import (
	"backend/internal/models"

	"github.com/jinzhu/gorm"
)

type SyncRepository interface {
	CreateSyncData(data *models.SyncData) error
	UpsertSyncData(data *models.SyncData) error
	DeleteSyncData(userID uint) error
	GetSyncData(userID uint) (*models.SyncData, error)
}

type syncRepository struct {
	db *gorm.DB
}

func NewSyncRepository(db *gorm.DB) SyncRepository {
	return &syncRepository{db}
}

func (r *syncRepository) CreateSyncData(data *models.SyncData) error {
	return r.db.Create(data).Error
}

func (r *syncRepository) UpsertSyncData(data *models.SyncData) error {
	var existingData models.SyncData
	err := r.db.Where("user_id = ?", data.UserID).First(&existingData).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return r.db.Create(data).Error
		}
		return err
	}
	existingData.ContentType = data.ContentType
	existingData.Age = data.Age
	existingData.QueryText = data.QueryText
	existingData.FeelingLevel = data.FeelingLevel
	existingData.ImageURL = data.ImageURL
	existingData.DocumentURL = data.DocumentURL

	return r.db.Save(&existingData).Error
}

func (r *syncRepository) DeleteSyncData(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.SyncData{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *syncRepository) GetSyncData(userID uint) (*models.SyncData, error) {
	var data models.SyncData
	err := r.db.Where("user_id = ?", userID).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
