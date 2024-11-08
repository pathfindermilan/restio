package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"sync"

	"github.com/jinzhu/gorm"
)

type SyncFeedbackService interface {
	UpsertSyncFeedback(data *models.SyncFeedback) error
	DeleteSyncFeedback(userID uint) error
	GetSyncFeedback(userID uint) (*models.SyncFeedback, error)
	ResetSyncFeedback(userID uint) error
}

type syncFeedback struct {
	repo repositories.SyncFeedbackRepository
	mu   sync.Mutex
}

func NewSyncFeedbackService(
	repo repositories.SyncFeedbackRepository,
) SyncFeedbackService {
	return &syncFeedback{
		repo: repo,
	}
}

func (s *syncFeedback) UpsertSyncFeedback(data *models.SyncFeedback) error {
	return s.repo.CreateOrUpdateSyncFeedback(data)
}

func (s *syncFeedback) DeleteSyncFeedback(userID uint) error {
	return s.repo.DeleteSyncFeedbackByUserID(userID)
}

func (s *syncFeedback) GetSyncFeedback(userID uint) (*models.SyncFeedback, error) {
	return s.repo.GetSyncFeedbackByUserID(userID)
}

func (s *syncFeedback) ResetSyncFeedback(userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	syncFeedback, err := s.GetSyncFeedback(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	updated := false
	if syncFeedback.HumanFeedback != "" || syncFeedback.HumanReaction != "" {
		syncFeedback.HumanFeedback = ""
		syncFeedback.HumanReaction = ""
		syncFeedback.NewContentType = ""
		updated = true
	}

	if updated {
		return s.UpsertSyncFeedback(syncFeedback)
	}

	return nil
}
