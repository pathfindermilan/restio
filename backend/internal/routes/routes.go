package routes

import (
	"backend/internal/auth"
	"backend/internal/controllers"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authService services.AuthService, aiService services.AIService, jwtService auth.JWTService) {
	authController := controllers.NewAuthController(authService)
	aiController := controllers.NewAIController(aiService)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)

	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware(jwtService))
	{
		protected.POST("/process", aiController.ProcessData)
	}
}
