package controllers

import (
	"net/http"

	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AIController struct {
	aiService services.AIService
}

func NewAIController(aiService services.AIService) *AIController {
	return &AIController{aiService}
}

func (ctrl *AIController) ProcessData(c *gin.Context) {
	var input struct {
		Data string `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := ctrl.aiService.ProcessData(input.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}
