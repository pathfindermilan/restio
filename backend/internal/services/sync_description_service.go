package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type SyncDescriptionService interface {
    CreateOrUpdateSyncDescription(data *models.SyncDescription) error
    GetSyncDescription(userID uint) (*models.SyncDescription, error)
    DeleteSyncDescription(userID uint) error
}

type syncDescriptionService struct {
    syncDescriptionRepo repositories.SyncDescriptionRepository
}

func NewSyncDescriptionService(syncDescriptionRepo repositories.SyncDescriptionRepository) SyncDescriptionService {
    return &syncDescriptionService{syncDescriptionRepo}
}

func (s *syncDescriptionService) CreateOrUpdateSyncDescription(data *models.SyncDescription) error {
    return s.syncDescriptionRepo.CreateOrUpdateSyncDescription(data)
}

func (s *syncDescriptionService) GetSyncDescription(userID uint) (*models.SyncDescription, error) {
    return s.syncDescriptionRepo.GetSyncDescriptionByUserID(userID)
}

func (s *syncDescriptionService) DeleteSyncDescription(userID uint) error {
    return s.syncDescriptionRepo.DeleteSyncDescription(userID)
}
