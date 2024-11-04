package main

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/routes"
	"backend/internal/services"

	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()

	dbURI := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + " dbname=" + cfg.DBName + " password=" + cfg.DBPassword + " sslmode=disable"
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	db.AutoMigrate(&models.User{}, &models.SyncData{}, &models.SyncDescription{})

	userRepo := repositories.NewUserRepository(db)
	syncRepo := repositories.NewSyncRepository(db)
	syncDescriptionRepo := repositories.NewSyncDescriptionRepository(db)

	jwtService := auth.NewJWTService(cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, jwtService)
	syncDescriptionService := services.NewSyncDescriptionService(syncDescriptionRepo)
	syncService := services.NewSyncService(userRepo, syncRepo, syncDescriptionService, &cfg)

	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://restio.xyz"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	// router.Static("/uploads/images", "./uploads/images")
	// router.Static("/uploads/documents", "./uploads/documents")

	routes.SetupRoutes(router, authService, jwtService, syncService, syncDescriptionService)

	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
