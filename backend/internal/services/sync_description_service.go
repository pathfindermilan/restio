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
	"sync"

	"github.com/jinzhu/gorm"
)

type SyncDescriptionService interface {
	GetSyncDescription(userID uint) (*models.SyncDescription, error)
	UpsertSyncDescription(desc *models.SyncDescription) error
	DeleteSyncDescription(userID uint) error
	ProcessSyncDescription(userID uint, processImage bool, processDocument bool)
}

type syncDescriptionService struct {
	syncDescriptionRepo repositories.SyncDescriptionRepository
	syncService         SyncService
	config              *config.Config
	mu                  sync.Mutex
}

func NewSyncDescriptionService(
	syncDescriptionRepo repositories.SyncDescriptionRepository,
	syncService SyncService,
	config *config.Config,
) SyncDescriptionService {
	return &syncDescriptionService{
		syncService:         syncService,
		syncDescriptionRepo: syncDescriptionRepo,
		config:              config,
	}
}

func (s *syncDescriptionService) GetSyncDescription(userID uint) (*models.SyncDescription, error) {
	return s.syncDescriptionRepo.GetSyncDescriptionByUserID(userID)
}

func (s *syncDescriptionService) UpsertSyncDescription(desc *models.SyncDescription) error {
	return s.syncDescriptionRepo.CreateOrUpdateSyncDescription(desc)
}

func (s *syncDescriptionService) DeleteSyncDescription(userID uint) error {
	return s.syncDescriptionRepo.DeleteSyncDescriptionByUserID(userID)
}

func (s *syncDescriptionService) ProcessSyncDescription(userID uint, newImageUploaded bool, newDocumentUploaded bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	syncData, err := s.syncService.GetSyncData(userID)
	if err != nil {
		fmt.Println("Error fetching SyncData:", err)
		return
	}

	existingSyncDescription, err := s.syncDescriptionRepo.GetSyncDescriptionByUserID(userID)
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
		err = s.syncDescriptionRepo.CreateOrUpdateSyncDescription(syncDescription)
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

func (s *syncDescriptionService) processImageAndUpdateDescription(userID uint, imageURL string) {
	imageSummary, status := s.processImage(imageURL)
	desc, err := s.GetSyncDescription(userID)
	if err != nil {
		fmt.Println("Error fetching SyncDescription for image processing:", err)
		return
	}
	desc.ImageSummary = imageSummary
	desc.ImageStatus = status
	err = s.UpsertSyncDescription(desc)
	if err != nil {
		fmt.Println("Error updating SyncDescription after image processing:", err)
	}
}

func (s *syncDescriptionService) processDocumentAndUpdateDescription(userID uint, documentURL string) {
	documentSummary, status := s.processDocument(documentURL)
	desc, err := s.GetSyncDescription(userID)
	if err != nil {
		fmt.Println("Error fetching SyncDescription for document processing:", err)
		return
	}
	desc.DocumentSummary = documentSummary
	desc.DocumentStatus = status
	err = s.UpsertSyncDescription(desc)
	if err != nil {
		fmt.Println("Error updating SyncDescription after document processing:", err)
	}
}

func (s *syncDescriptionService) processImage(imagePath string) (string, models.SyncStatus) {
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

func (s *syncDescriptionService) processDocument(documentPath string) (string, models.SyncStatus) {
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
