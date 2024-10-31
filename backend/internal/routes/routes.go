package routes

import (
	"backend/internal/auth"
	"backend/internal/controllers"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authService services.AuthService, jwtService auth.JWTService) {
	authController := controllers.NewAuthController(authService)
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
	}
}
