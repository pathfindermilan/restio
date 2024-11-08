package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/jinzhu/gorm"
)

type AISummaryService interface {
	UpsertAISummary(data *models.AISummary) error
	DeleteAISummary(userID uint) error
	GetAISummary(userID uint) (*models.AISummary, error)
	ResetAriaResponse(userID uint) error
	ResetAllegroResponse(userID uint) error
	GenerateAISummaryAnswer(userID string) (string, models.SyncStatus)
}

type aiSummaryService struct {
	repo                   repositories.AISummaryRepository
	userRepo               repositories.UserRepository
	syncService            SyncService
	syncDescriptionService SyncDescriptionService
	syncFeedbackService    SyncFeedbackService
	config                 *config.Config
	mu                     sync.Mutex
}

func NewAISummaryService(
	repo repositories.AISummaryRepository,
	userRepo repositories.UserRepository,
	syncService SyncService,
	syncDescriptionService SyncDescriptionService,
	syncFeedbackService SyncFeedbackService,
	config *config.Config,
) AISummaryService {
	return &aiSummaryService{
		repo:                   repo,
		userRepo:               userRepo,
		syncService:            syncService,
		syncDescriptionService: syncDescriptionService,
		syncFeedbackService:    syncFeedbackService,
		config:                 config,
	}
}

func (s *aiSummaryService) UpsertAISummary(data *models.AISummary) error {
	return s.repo.CreateOrUpdateAISummary(data)
}

func (s *aiSummaryService) DeleteAISummary(userID uint) error {
	return s.repo.DeleteAISummaryByUserID(userID)
}

func (s *aiSummaryService) GetAISummary(userID uint) (*models.AISummary, error) {
	return s.repo.GetAISummaryByUserID(userID)
}

func (s *aiSummaryService) ResetAriaResponse(userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	aiSummary, err := s.GetAISummary(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	updated := false
	if aiSummary.AriaResponse != "" || aiSummary.AriaResponseStatus != models.StatusInProgress {
		aiSummary.AriaResponse = ""
		aiSummary.AriaResponseStatus = models.StatusInProgress
		updated = true
	}

	if updated {
		return s.UpsertAISummary(aiSummary)
	}

	return nil
}

func (s *aiSummaryService) ResetAllegroResponse(userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	aiSummary, err := s.GetAISummary(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	updated := false
	if aiSummary.AllegroResponse != "" || aiSummary.AllegroResponseStatus != models.StatusInProgress {
		aiSummary.AllegroResponse = ""
		aiSummary.AllegroResponseStatus = models.StatusInProgress
		updated = true
	}

	if updated {
		return s.UpsertAISummary(aiSummary)
	}

	return nil
}

func (s *aiSummaryService) GenerateAISummaryAnswer(user_id string) (string, models.SyncStatus) {
	uid, err := strconv.Atoi(user_id)
	if err != nil {
		return "", models.StatusErrored
	}
	userID := uint(uid)

	syncData, err := s.syncService.GetSyncData(userID)
	if err != nil {
		return "", models.StatusErrored
	}

	existingAISummary, err := s.repo.GetAISummaryByUserID(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		fmt.Println("Error fetching SyncDescription:", err)
		return "", models.StatusErrored
	}

	var aiSummary *models.AISummary
	if existingAISummary == nil {
		aiSummary = &models.AISummary{
			UserID:     userID,
			SyncDataID: syncData.ID,
		}
	} else {
		aiSummary = existingAISummary
	}

	syncDesc, err := s.syncDescriptionService.GetSyncDescription(userID)
	if err != nil {
		return "", models.StatusErrored
	}

	syncFeedback, err := s.syncFeedbackService.GetSyncFeedback(userID)
	if err != nil {
		return "", models.StatusErrored
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return "", models.StatusErrored
	}

	var previous_answer string
	if aiSummary.AllegroResponseStatus == models.StatusDone {
		previous_answer = aiSummary.AriaResponse
	} else {
		previous_answer = ""
	}

	payload := map[string]interface{}{
		"name":             user.Name,
		"age":              syncData.Age,
		"image_summary":    syncDesc.ImageSummary,
		"document_summary": syncDesc.DocumentSummary,
		"user_transcript":  syncData.QueryText,
		"content_type":     syncData.ContentType,
		"feeling_level":    syncData.FeelingLevel,
		"human_feedback":   syncFeedback.HumanFeedback,
		"human_reaction":   syncFeedback.HumanReaction,
		"previous_answer":  previous_answer,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", models.StatusErrored
	}

	req, err := http.NewRequest("POST", s.config.GenerateAnswerEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", models.StatusErrored
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", models.StatusErrored
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", models.StatusErrored
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("FastAPI responded with status: %d, body: %s\n", resp.StatusCode, string(body))
		return "", models.StatusErrored
	}

	var responseData struct {
		AI_Summary string `json:"ai_summary"`
	}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", models.StatusErrored
	}

	if syncData.ContentType != "video" {
		aiSummary.AriaResponse = responseData.AI_Summary
		aiSummary.AriaResponseStatus = models.StatusDone
	} else {
		aiSummary.AllegroResponse = responseData.AI_Summary
		aiSummary.AllegroResponseStatus = models.StatusDone
	}

	err = s.repo.CreateOrUpdateAISummary(aiSummary)
	if err != nil {
		return "", models.StatusErrored
	}

	if syncData.ContentType != "video" {
		return existingAISummary.AriaResponse, models.StatusDone
	} else {
		return existingAISummary.AllegroResponse, models.StatusDone
	}
}
