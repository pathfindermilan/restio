package repositories

import (
	"backend/internal/models"

	"github.com/jinzhu/gorm"
)

type AISummaryRepository interface {
	GetAISummaryByUserID(userID uint) (*models.AISummary, error)
	CreateOrUpdateAISummary(desc *models.AISummary) error
	DeleteAISummaryByUserID(userID uint) error
}

type aiSummaryRepository struct {
	db *gorm.DB
}

func NewAISummaryRepository(db *gorm.DB) AISummaryRepository {
	return &aiSummaryRepository{db}
}

func (r *aiSummaryRepository) GetAISummaryByUserID(userID uint) (*models.AISummary, error) {
	var desc models.AISummary
	err := r.db.Where("user_id = ?", userID).First(&desc).Error
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

func (r *aiSummaryRepository) CreateOrUpdateAISummary(desc *models.AISummary) error {
	var existing models.AISummary
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

func (r *aiSummaryRepository) DeleteAISummaryByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.AISummary{}).Error
}
