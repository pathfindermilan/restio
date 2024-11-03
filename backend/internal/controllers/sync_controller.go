package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/microcosm-cc/bluemonday"
)

type SyncController struct {
	syncService            services.SyncService
	syncDescriptionService services.SyncDescriptionService
}

func NewSyncController(
	syncService services.SyncService,
	syncDescriptionService services.SyncDescriptionService,
) *SyncController {
	return &SyncController{
		syncService:            syncService,
		syncDescriptionService: syncDescriptionService,
	}
}

func (ctrl *SyncController) SyncData(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDInterface.(uint)

	contentType := c.PostForm("content_type")
	ageStr := c.PostForm("age")
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age"})
		return
	}
	queryText := c.PostForm("query_text")
	feelingLevel := c.PostForm("feeling_level")

	allowedContentTypes := map[string]bool{
		"joke":     true,
		"speech":   true,
		"exercise": true,
		"video":    true,
	}
	if !allowedContentTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content_type"})
		return
	}

	allowedFeelings := map[string]bool{
		"happy":     true,
		"sad":       true,
		"angry":     true,
		"confused":  true,
		"confident": true,
		"tired":     true,
	}
	if !allowedFeelings[feelingLevel] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feeling_level"})
		return
	}

	imageFile, _ := c.FormFile("image")
	documentFile, _ := c.FormFile("document")

	var imageURL, documentURL string
	var uploadImage bool = false
	var uploadDocument bool = false

	allowedImageTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	allowedDocumentTypes := map[string]bool{
		"application/pdf":    true,
		"application/msword": true,
	}

	policy := bluemonday.NewPolicy()

	existingSyncData, err := ctrl.syncService.GetSyncData(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve existing sync data"})
		return
	}

	if imageFile != nil {
		mimeType := imageFile.Header.Get("Content-Type")
		if !allowedImageTypes[mimeType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file type"})
			return
		}

		cleanImageFilename := policy.Sanitize(imageFile.Filename)

		uploadImage = true
		if existingSyncData != nil && existingSyncData.ImageURL != "" {
			existingImageName := filepath.Base(existingSyncData.ImageURL)
			if strings.EqualFold(cleanImageFilename, existingImageName) {
				uploadImage = false
			}
		}

		if uploadImage {
			if existingSyncData != nil && existingSyncData.ImageURL != "" {
				if err := os.Remove(existingSyncData.ImageURL); err != nil && !os.IsNotExist(err) {
					c.Writer.WriteString("Warning: Failed to delete old image file\n")
				}
			}

			newImageFilename := strconv.Itoa(int(userID)) + "-" + imageFile.Filename
			imagePath := filepath.Join("uploads/images", newImageFilename)

			if err := os.MkdirAll(filepath.Dir(imagePath), os.ModePerm); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image directory"})
				return
			}

			if err := c.SaveUploadedFile(imageFile, imagePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
				return
			}
			imageURL = imagePath
		} else {
			imageURL = existingSyncData.ImageURL
		}
	} else {
		if existingSyncData != nil {
			imageURL = existingSyncData.ImageURL
		}
	}

	if documentFile != nil {
		mimeType := documentFile.Header.Get("Content-Type")
		if !allowedDocumentTypes[mimeType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document file type"})
			return
		}

		cleanDocumentFilename := policy.Sanitize(documentFile.Filename)

		uploadDocument = true
		if existingSyncData != nil && existingSyncData.DocumentURL != "" {
			existingDocumentName := filepath.Base(existingSyncData.DocumentURL)
			if strings.EqualFold(cleanDocumentFilename, existingDocumentName) {
				uploadDocument = false
			}
		}

		if uploadDocument {
			if existingSyncData != nil && existingSyncData.DocumentURL != "" {
				if err := os.Remove(existingSyncData.DocumentURL); err != nil && !os.IsNotExist(err) {
					c.Writer.WriteString("Warning: Failed to delete old document file\n")
				}
			}

			newDocumentFilename := strconv.Itoa(int(userID)) + "-" + documentFile.Filename
			documentPath := filepath.Join("uploads/documents", newDocumentFilename)

			if err := os.MkdirAll(filepath.Dir(documentPath), os.ModePerm); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create documents directory"})
				return
			}

			if err := c.SaveUploadedFile(documentFile, documentPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
				return
			}
			documentURL = documentPath
		} else {
			documentURL = existingSyncData.DocumentURL
		}
	} else {
		if existingSyncData != nil {
			documentURL = existingSyncData.DocumentURL
		}
	}

	syncData := models.SyncData{
		UserID:       userID,
		ContentType:  contentType,
		Age:          age,
		QueryText:    queryText,
		FeelingLevel: feelingLevel,
		ImageURL:     imageURL,
		DocumentURL:  documentURL,
	}

	newImageUploaded := (imageFile != nil && uploadImage)
	newDocumentUploaded := (documentFile != nil && uploadDocument)

	err = ctrl.syncService.UpsertSyncData(&syncData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save sync data"})
		return
	}

	go ctrl.syncService.ProcessDescriptions(userID, newImageUploaded, newDocumentUploaded)

	c.JSON(http.StatusOK, gin.H{"message": "Data synced successfully"})
}

func (ctrl *SyncController) SyncReset(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDInterface.(uint)

	syncData, err := ctrl.syncService.GetSyncData(userID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No sync data found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sync data"})
		return
	}

	if syncData.ImageURL != "" {
		if err := os.Remove(syncData.ImageURL); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image file"})
			return
		}
	}

	if syncData.DocumentURL != "" {
		if err := os.Remove(syncData.DocumentURL); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document file"})
			return
		}
	}

	deleteErr := ctrl.syncService.DeleteSyncData(userID)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sync data"})
		return
	}

	err = ctrl.syncDescriptionService.DeleteSyncDescription(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sync description data"})
		return
	}

	_, err = ctrl.syncService.GetSyncData(userID)
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sync data record still exists after deletion"})
		return
	} else if !gorm.IsRecordNotFoundError(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error after deletion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync data reset successfully"})
}

func (ctrl *SyncController) ServeImage(c *gin.Context) {
	filename := c.Param("filename")
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDInterface.(uint)

	filename = filepath.Base(filename)

	split := strings.SplitN(filename, "-", 2)

	fileUserIDStr := split[0]
	fileUserID, err := strconv.Atoi(fileUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in filename"})
		return
	}

	if uint(fileUserID) != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this file"})
		return
	}

	filePath := filepath.Join("uploads/images", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}

func (ctrl *SyncController) ServeDocument(c *gin.Context) {
	filename := c.Param("filename")
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDInterface.(uint)

	filename = filepath.Base(filename)

	split := strings.SplitN(filename, "-", 2)

	fileUserIDStr := split[0]
	fileUserID, err := strconv.Atoi(fileUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in filename"})
		return
	}

	if uint(fileUserID) != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this file"})
		return
	}

	filePath := filepath.Join("uploads/documents", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}
