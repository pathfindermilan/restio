package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SyncController struct {
    syncService services.SyncService
}

func NewSyncController(syncService services.SyncService) *SyncController {
    return &SyncController{syncService}
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
        "joke": true,
        "speech": true,
        "exercise": true,
        "video": true,
    }
    if !allowedContentTypes[contentType] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content_type"})
        return
    }

    allowedFeelings := map[string]bool{
        "happy": true,
        "sad": true,
        "angry": true,
        "confused": true,
        "confident": true,
        "tired": true,
    }
    if !allowedFeelings[feelingLevel] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feeling_level"})
        return
    }

    imageFile, _ := c.FormFile("image")
    documentFile, _ := c.FormFile("document")

    var imageURL, documentURL string

    if imageFile != nil {
        imagePath := filepath.Join("uploads/images", imageFile.Filename)
        if err := c.SaveUploadedFile(imageFile, imagePath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
            return
        }
        imageURL = imagePath
    }

    if documentFile != nil {
        documentPath := filepath.Join("uploads/documents", documentFile.Filename)
        if err := c.SaveUploadedFile(documentFile, documentPath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
            return
        }
        documentURL = documentPath
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

    err = ctrl.syncService.CreateSyncData(&syncData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Data synced successfully"})
}
