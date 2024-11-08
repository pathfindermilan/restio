package routes

import (
	"backend/internal/auth"
	"backend/internal/controllers"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authService services.AuthService,
	jwtService auth.JWTService,
	syncService services.SyncService,
	syncDescriptionService services.SyncDescriptionService,
) {
	authController := controllers.NewAuthController(authService)
	syncController := controllers.NewSyncController(syncService, syncDescriptionService)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.POST("/verify-email", authController.VerifyEmail)

	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware(jwtService))
	{
		protected.POST("/logout", authController.Logout)
		protected.DELETE("/delete", authController.DeleteUser)
		protected.PUT("/update", authController.UpdateUser)
		protected.GET("/profile", authController.GetProfile)

		protected.POST("/sync", syncController.SyncData)
		protected.DELETE("/sync-reset", syncController.SyncReset)

		protected.GET("/uploads/images/:filename", syncController.ServeImage)
		protected.GET("/uploads/documents/:filename", syncController.ServeDocument)
	}
}
