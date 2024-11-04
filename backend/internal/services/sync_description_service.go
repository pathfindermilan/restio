package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type SyncDescriptionService interface {
	GetSyncDescriptionByUserID(userID uint) (*models.SyncDescription, error)
	CreateOrUpdateSyncDescription(desc *models.SyncDescription) error
	DeleteSyncDescription(userID uint) error
}

type syncDescriptionService struct {
	repo repositories.SyncDescriptionRepository
}

func NewSyncDescriptionService(repo repositories.SyncDescriptionRepository) SyncDescriptionService {
	return &syncDescriptionService{repo}
}

func (s *syncDescriptionService) GetSyncDescriptionByUserID(userID uint) (*models.SyncDescription, error) {
	return s.repo.GetSyncDescriptionByUserID(userID)
}

func (s *syncDescriptionService) CreateOrUpdateSyncDescription(desc *models.SyncDescription) error {
	return s.repo.CreateOrUpdateSyncDescription(desc)
}

func (s *syncDescriptionService) DeleteSyncDescription(userID uint) error {
	return s.repo.DeleteSyncDescription(userID)
}
