package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type SyncService interface {
	UpsertSyncData(data *models.SyncData) error
	DeleteSyncData(userID uint) error
	GetSyncData(userID uint) (*models.SyncData, error)
}

type syncService struct {
	repo repositories.SyncRepository
}

func NewSyncService(
	repo repositories.SyncRepository,
) SyncService {
	return &syncService{
		repo: repo,
	}
}

func (s *syncService) UpsertSyncData(data *models.SyncData) error {
	return s.repo.CreateOrUpdateSyncData(data)
}

func (s *syncService) DeleteSyncData(userID uint) error {
	return s.repo.DeleteSyncDataByUserID(userID)
}

func (s *syncService) GetSyncData(userID uint) (*models.SyncData, error) {
	return s.repo.GetSyncDataByUserID(userID)
}
