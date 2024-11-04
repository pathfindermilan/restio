package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repositories"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"fmt"

	"github.com/jinzhu/gorm"
)

type SyncService interface {
	CreateSyncData(data *models.SyncData) error
	UpsertSyncData(data *models.SyncData) error
	DeleteSyncData(userID uint) error
	GetSyncData(userID uint) (*models.SyncData, error)
	ProcessDescriptions(userID uint, processImage bool, processDocument bool)
	DeleteSyncDescription(userID uint) error
}

type syncService struct {
	syncRepo               repositories.SyncRepository
	syncDescriptionService SyncDescriptionService
	config                 *config.Config
	mu                     sync.Mutex
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

func (s *syncService) DeleteSyncDescription(userID uint) error {
	return s.syncDescriptionService.DeleteSyncDescription(userID)
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
