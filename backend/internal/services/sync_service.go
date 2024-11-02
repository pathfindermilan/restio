package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type SyncService interface {
    CreateSyncData(data *models.SyncData) error
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
