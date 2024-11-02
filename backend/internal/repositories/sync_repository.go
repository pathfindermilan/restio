package repositories

import (
	"backend/internal/models"

	"github.com/jinzhu/gorm"
)

type SyncRepository interface {
    CreateSyncData(data *models.SyncData) error
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
