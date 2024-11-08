package repositories

import (
	"backend/internal/models"

	"github.com/jinzhu/gorm"
)

type SyncRepository interface {
	CreateOrUpdateSyncData(data *models.SyncData) error
	DeleteSyncDataByUserID(userID uint) error
	GetSyncDataByUserID(userID uint) (*models.SyncData, error)
}

type syncRepository struct {
	db *gorm.DB
}

func NewSyncRepository(db *gorm.DB) SyncRepository {
	return &syncRepository{db}
}

func (r *syncRepository) CreateOrUpdateSyncData(data *models.SyncData) error {
	var existing models.SyncData
	err := r.db.Where("user_id = ?", data.UserID).First(&existing).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if existing.ID != 0 {
		data.ID = existing.ID
		return r.db.Save(data).Error
	}
	return r.db.Create(data).Error
}

func (r *syncRepository) DeleteSyncDataByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.SyncData{}).Error
}

func (r *syncRepository) GetSyncDataByUserID(userID uint) (*models.SyncData, error) {
	var data models.SyncData
	err := r.db.Where("user_id = ?", userID).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
