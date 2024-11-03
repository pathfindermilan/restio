package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type SyncService interface {
	CreateSyncData(data *models.SyncData) error
	UpsertSyncData(data *models.SyncData) error
	DeleteSyncData(userID uint) error
	GetSyncData(userID uint) (*models.SyncData, error)
}

type syncService struct {
	syncRepo repositories.SyncRepository
}

func NewSyncService(syncRepo repositories.SyncRepository) SyncService {
	return &syncService{syncRepo}
}

func (s *syncService) CreateSyncData(data *models.SyncData) error {
	return s.syncRepo.CreateSyncData(data)
}

func (s *syncService) UpsertSyncData(data *models.SyncData) error {
	return s.syncRepo.UpsertSyncData(data)
}

func (s *syncService) DeleteSyncData(userID uint) error {
	return s.syncRepo.DeleteSyncData(userID)
}

func (s *syncService) GetSyncData(userID uint) (*models.SyncData, error) {
	return s.syncRepo.GetSyncData(userID)
}
