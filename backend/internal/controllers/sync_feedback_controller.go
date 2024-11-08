package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SyncFeedbackController struct {
	syncService         services.SyncService
	syncFeedbackService services.SyncFeedbackService
	aiSummaryService    services.AISummaryService
}

func NewSyncFeedbackController(
	syncService services.SyncService,
	syncFeedbackService services.SyncFeedbackService,
	aiSummaryService services.AISummaryService,
) *SyncFeedbackController {
	return &SyncFeedbackController{
		syncService:         syncService,
		syncFeedbackService: syncFeedbackService,
		aiSummaryService:    aiSummaryService,
	}
}

func (ctrl *SyncFeedbackController) SyncFeedackData(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDInterface.(uint)

	contentType := c.PostForm("content_type")

	queryText := c.PostForm("query_text")

	feelingLevelStr := c.PostForm("feeling_level")
	feelingLevelInt, err := strconv.Atoi(feelingLevelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feeling_level"})
		return
	}

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

	feelingLevelMap := map[int]string{
		1:  "very sad",
		2:  "sad",
		3:  "slightly down",
		4:  "neutral low",
		5:  "okay",
		6:  "slightly good",
		7:  "good",
		8:  "very good",
		9:  "great",
		10: "excellent",
	}

	feelingLevel, exists := feelingLevelMap[feelingLevelInt]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feeling_level"})
		return
	}

	existingSyncData, err := ctrl.syncService.GetSyncData(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve existing sync data"})
		return
	}

	existingSyncFeedback, err := ctrl.syncFeedbackService.GetSyncFeedback(userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve existing sync data"})
		return
	}

	syncFeedback := &models.SyncFeedback{
		UserID:         userID,
		SyncDataID:     existingSyncData.ID,
		NewContentType: contentType,
		HumanFeedback:  queryText,
		HumanReaction:  feelingLevel,
		QuestionsRelly: existingSyncFeedback.QuestionsRelly + 1,
	}

	err = ctrl.syncFeedbackService.UpsertSyncFeedback(syncFeedback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save sync data"})
		return
	}

	// call the ai model after updating the feedback
	//
	c.JSON(http.StatusOK, gin.H{"message": "Feedback synced successfully"})
}
