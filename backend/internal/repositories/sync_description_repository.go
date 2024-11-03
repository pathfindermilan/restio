package repositories

import (
	"backend/internal/models"

	"github.com/jinzhu/gorm"
)

type SyncDescriptionRepository interface {
    CreateOrUpdateSyncDescription(data *models.SyncDescription) error
    GetSyncDescriptionByUserID(userID uint) (*models.SyncDescription, error)
    DeleteSyncDescription(userID uint) error
}

type syncDescriptionRepository struct {
    db *gorm.DB
}

func NewSyncDescriptionRepository(db *gorm.DB) SyncDescriptionRepository {
    return &syncDescriptionRepository{db}
}

func (r *syncDescriptionRepository) CreateOrUpdateSyncDescription(data *models.SyncDescription) error {
    var existing models.SyncDescription
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

func (r *syncDescriptionRepository) GetSyncDescriptionByUserID(userID uint) (*models.SyncDescription, error) {
    var data models.SyncDescription
    err := r.db.Where("user_id = ?", userID).First(&data).Error
    return &data, err
}

func (r *syncDescriptionRepository) DeleteSyncDescription(userID uint) error {
    return r.db.Where("user_id = ?", userID).Delete(&models.SyncDescription{}).Error
}
