package repositories

import (
	"backend/internal/models"

	"github.com/jinzhu/gorm"
)

type SyncFeedbackRepository interface {
	GetSyncFeedbackByUserID(userID uint) (*models.SyncFeedback, error)
	CreateOrUpdateSyncFeedback(desc *models.SyncFeedback) error
	DeleteSyncFeedbackByUserID(userID uint) error
}

type syncFeedbackRepository struct {
	db *gorm.DB
}

func NewSyncFeedbackRepository(db *gorm.DB) SyncFeedbackRepository {
	return &syncFeedbackRepository{db}
}

func (r *syncFeedbackRepository) GetSyncFeedbackByUserID(userID uint) (*models.SyncFeedback, error) {
	var desc models.SyncFeedback
	err := r.db.Where("user_id = ?", userID).First(&desc).Error
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

func (r *syncFeedbackRepository) CreateOrUpdateSyncFeedback(desc *models.SyncFeedback) error {
	var existing models.SyncFeedback
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

func (r *syncFeedbackRepository) DeleteSyncFeedbackByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.SyncFeedback{}).Error
}
