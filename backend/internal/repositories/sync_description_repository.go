package repositories

import (
	"backend/internal/models"

	"github.com/jinzhu/gorm"
)

type SyncDescriptionRepository interface {
	GetSyncDescriptionByUserID(userID uint) (*models.SyncDescription, error)
	CreateOrUpdateSyncDescription(desc *models.SyncDescription) error
	DeleteSyncDescription(userID uint) error
}

type syncDescriptionRepository struct {
	db *gorm.DB
}

func NewSyncDescriptionRepository(db *gorm.DB) SyncDescriptionRepository {
	return &syncDescriptionRepository{db}
}

func (r *syncDescriptionRepository) GetSyncDescriptionByUserID(userID uint) (*models.SyncDescription, error) {
	var desc models.SyncDescription
	err := r.db.Where("user_id = ?", userID).First(&desc).Error
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

func (r *syncDescriptionRepository) CreateOrUpdateSyncDescription(desc *models.SyncDescription) error {
	var existing models.SyncDescription
	err := r.db.Where("user_id = ?", desc.UserID).First(&existing).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if existing.ID != 0 {
		desc.ID = existing.ID
		return r.db.Save(desc).Error
	}
	return r.db.Create(desc).Error
}

func (r *syncDescriptionRepository) DeleteSyncDescription(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.SyncDescription{}).Error
}
