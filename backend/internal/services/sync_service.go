package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/jinzhu/gorm"
)

type SyncService interface {
	CreateSyncData(data *models.SyncData) error
	UpsertSyncData(data *models.SyncData) error
	DeleteSyncData(userID uint) error
	GetSyncData(userID uint) (*models.SyncData, error)
	ProcessDescriptions(userID uint, processImage bool, processDocument bool)
	DeleteSyncDescription(userID uint) error
	ResetAIResponse(userID uint) error
	GenerateAnswer(userID string) (string, models.SyncStatus)
}

type syncService struct {
	userRepo               repositories.UserRepository
	syncRepo               repositories.SyncRepository
	syncDescriptionService SyncDescriptionService
	config                 *config.Config
	mu                     sync.Mutex
}

func NewSyncService(
	userRepo repositories.UserRepository,
	syncRepo repositories.SyncRepository,
	syncDescriptionService SyncDescriptionService,
	config *config.Config,
) SyncService {
	return &syncService{
		userRepo:               userRepo,
		syncRepo:               syncRepo,
		syncDescriptionService: syncDescriptionService,
		config:                 config,
	}
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

func (s *syncService) DeleteSyncDescription(userID uint) error {
	return s.syncDescriptionService.DeleteSyncDescription(userID)
}

func (s *syncService) ResetAIResponse(userID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	syncDesc, err := s.syncDescriptionService.GetSyncDescriptionByUserID(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	if syncDesc == nil {
		syncDesc = &models.SyncDescription{
			UserID:     userID,
			AiSummary:  "",
			AiStatus:   models.StatusInProgress,
			SyncDataID: 0, // Assuming 0 or appropriate value
		}
		return s.syncDescriptionService.CreateOrUpdateSyncDescription(syncDesc)
	}

	updated := false
	if syncDesc.AiSummary != "" || syncDesc.AiStatus != models.StatusInProgress {
		syncDesc.AiSummary = ""
		syncDesc.AiStatus = models.StatusInProgress
		updated = true
	}

	if updated {
		return s.syncDescriptionService.CreateOrUpdateSyncDescription(syncDesc)
	}

	return nil
}

func (s *syncService) ProcessDescriptions(userID uint, newImageUploaded bool, newDocumentUploaded bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	syncData, err := s.GetSyncData(userID)
	if err != nil {
		fmt.Println("Error fetching SyncData:", err)
		return
	}

	existingSyncDescription, err := s.syncDescriptionService.GetSyncDescriptionByUserID(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		fmt.Println("Error fetching SyncDescription:", err)
		return
	}

	var syncDescription *models.SyncDescription
	if existingSyncDescription == nil {
		syncDescription = &models.SyncDescription{
			UserID:     userID,
			SyncDataID: syncData.ID,
		}
	} else {
		syncDescription = existingSyncDescription
	}

	needToUpdate := false

	if newImageUploaded {
		if syncDescription.ImageURL != syncData.ImageURL {
			syncDescription.ImageURL = syncData.ImageURL
			syncDescription.ImageStatus = models.StatusInProgress
			syncDescription.ImageSummary = ""
			needToUpdate = true
		}
	} else {
		if syncData.ImageURL == "" && syncDescription.ImageStatus != models.StatusNotFound {
			syncDescription.ImageStatus = models.StatusNotFound
			needToUpdate = true
		}
	}

	if newDocumentUploaded {
		if syncDescription.DocumentURL != syncData.DocumentURL {
			syncDescription.DocumentURL = syncData.DocumentURL
			syncDescription.DocumentStatus = models.StatusInProgress
			syncDescription.DocumentSummary = ""
			needToUpdate = true
		}
	} else {
		if syncData.DocumentURL == "" && syncDescription.DocumentStatus != models.StatusNotFound {
			syncDescription.DocumentStatus = models.StatusNotFound
			needToUpdate = true
		}
	}

	if existingSyncDescription == nil && !newImageUploaded && !newDocumentUploaded {
		if syncData.ImageURL == "" {
			syncDescription.ImageStatus = models.StatusNotFound
		}
		if syncData.DocumentURL == "" {
			syncDescription.DocumentStatus = models.StatusNotFound
		}
		needToUpdate = true
	}

	if existingSyncDescription == nil || needToUpdate {
		err = s.syncDescriptionService.CreateOrUpdateSyncDescription(syncDescription)
		if err != nil {
			fmt.Println("Error saving SyncDescription:", err)
			return
		}
	}

	if newImageUploaded && syncDescription.ImageStatus == models.StatusInProgress {
		go s.processImageAndUpdateDescription(userID, syncData.ImageURL)
	}

	if newDocumentUploaded && syncDescription.DocumentStatus == models.StatusInProgress {
		go s.processDocumentAndUpdateDescription(userID, syncData.DocumentURL)
	}

	go func(uid uint) {
		answer, status := s.GenerateAnswer(strconv.Itoa(int(uid)))
		if status != models.StatusDone {
			fmt.Printf("Failed to generate answer for user %d\n", uid)
			return
		}
		fmt.Printf("Generated answer for user %d: %s\n", uid, answer)
	}(userID)
}

func (s *syncService) processImageAndUpdateDescription(userID uint, imageURL string) {
	imageSummary, status := s.processImage(imageURL)
	desc, err := s.syncDescriptionService.GetSyncDescriptionByUserID(userID)
	if err != nil {
		fmt.Println("Error fetching SyncDescription for image processing:", err)
		return
	}
	desc.ImageSummary = imageSummary
	desc.ImageStatus = status
	err = s.syncDescriptionService.CreateOrUpdateSyncDescription(desc)
	if err != nil {
		fmt.Println("Error updating SyncDescription after image processing:", err)
	}
}

func (s *syncService) processDocumentAndUpdateDescription(userID uint, documentURL string) {
	documentSummary, status := s.processDocument(documentURL)
	desc, err := s.syncDescriptionService.GetSyncDescriptionByUserID(userID)
	if err != nil {
		fmt.Println("Error fetching SyncDescription for document processing:", err)
		return
	}
	desc.DocumentSummary = documentSummary
	desc.DocumentStatus = status
	err = s.syncDescriptionService.CreateOrUpdateSyncDescription(desc)
	if err != nil {
		fmt.Println("Error updating SyncDescription after document processing:", err)
	}
}

func (s *syncService) processImage(imagePath string) (string, models.SyncStatus) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", models.StatusErrored
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return "", models.StatusErrored
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", models.StatusErrored
	}
	writer.Close()

	req, err := http.NewRequest("POST", s.config.DescribeImageEndpoint, body)
	if err != nil {
		return "", models.StatusErrored
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", models.StatusErrored
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", models.StatusErrored
	}

	var respData struct {
		ImageSummary string `json:"image_summary"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", models.StatusErrored
	}

	return respData.ImageSummary, models.StatusDone
}

func (s *syncService) processDocument(documentPath string) (string, models.SyncStatus) {
	file, err := os.Open(documentPath)
	if err != nil {
		return "", models.StatusErrored
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("document", filepath.Base(documentPath))
	if err != nil {
		return "", models.StatusErrored
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", models.StatusErrored
	}
	writer.Close()

	req, err := http.NewRequest("POST", s.config.DescribeDocumentEndpoint, body)
	if err != nil {
		return "", models.StatusErrored
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", models.StatusErrored
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", models.StatusErrored
	}

	var respData struct {
		DocumentSummary string `json:"document_summary"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", models.StatusErrored
	}

	return respData.DocumentSummary, models.StatusDone
}

func (s *syncService) GenerateAnswer(user_id string) (string, models.SyncStatus) {
	uid, err := strconv.Atoi(user_id)
	if err != nil {
		return "", models.StatusErrored
	}
	userID := uint(uid)

	syncData, err := s.GetSyncData(userID)
	if err != nil {
		return "", models.StatusErrored
	}

	syncDesc, err := s.syncDescriptionService.GetSyncDescriptionByUserID(userID)
	if err != nil {
		return "", models.StatusErrored
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return "", models.StatusErrored
	}
	name := user.Name

	payload := map[string]interface{}{
		"name":             name,
		"age":              syncData.Age,
		"image_summary":    syncDesc.ImageSummary,
		"document_summary": syncDesc.DocumentSummary,
		"user_transcript":  syncData.QueryText,
		"content_type":     syncData.ContentType,
		"feeling_level":    syncData.FeelingLevel,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", models.StatusErrored
	}

	req, err := http.NewRequest("POST", "https://restio.site/generate-answer", bytes.NewBuffer(payloadBytes))
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

	syncDesc.AiSummary = responseData.AI_Summary
	syncDesc.AiStatus = models.StatusDone

	err = s.syncDescriptionService.CreateOrUpdateSyncDescription(syncDesc)
	if err != nil {
		return "", models.StatusErrored
	}

	return syncDesc.AiSummary, models.StatusDone
}
