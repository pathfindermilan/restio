package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repositories"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type SyncService interface {
	CreateSyncData(data *models.SyncData) error
	UpsertSyncData(data *models.SyncData) error
	DeleteSyncData(userID uint) error
	GetSyncData(userID uint) (*models.SyncData, error)
	ProcessDescriptions(userID uint, processImage bool, processDocument bool)
}

type syncService struct {
	syncRepo               repositories.SyncRepository
	syncDescriptionService SyncDescriptionService
	config                 *config.Config
}

func NewSyncService(
	syncRepo repositories.SyncRepository,
	syncDescriptionService SyncDescriptionService,
	config *config.Config,
) SyncService {
	return &syncService{
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

func (s *syncService) ProcessDescriptions(userID uint, processImageFlag bool, processDocumentFlag bool) {
	syncData, err := s.GetSyncData(userID)
	if err != nil {
		fmt.Println("Error fetching sync data:", err)
		return
	}

	syncDescription := &models.SyncDescription{
		UserID:     userID,
		SyncDataID: syncData.ID,
	}

	if processImageFlag && syncData.ImageURL != "" {
		syncDescription.ImageStatus = models.StatusInProgress
	} else {
		syncDescription.ImageStatus = models.StatusNotFound
	}

	if processDocumentFlag && syncData.DocumentURL != "" {
		syncDescription.DocumentStatus = models.StatusInProgress
	} else {
		syncDescription.DocumentStatus = models.StatusNotFound
	}

	if syncDescription.ImageStatus == models.StatusInProgress {
		imageSummary, status := s.processImage(syncData.ImageURL)
		syncDescription.ImageSummary = imageSummary
		syncDescription.ImageStatus = status
	}

	if syncDescription.DocumentStatus == models.StatusInProgress {
		documentSummary, status := s.processDocument(syncData.DocumentURL)
		syncDescription.DocumentSummary = documentSummary
		syncDescription.DocumentStatus = status
	}

	err = s.syncDescriptionService.CreateOrUpdateSyncDescription(syncDescription)
	if err != nil {
		fmt.Println("Error saving sync description:", err)
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
