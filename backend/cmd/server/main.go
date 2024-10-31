package main

import (
	"backend/internal/models"

	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/repositories"
	"backend/internal/routes"
	"backend/internal/services"

	"log"

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

	db.AutoMigrate(&models.User{})

	userRepo := repositories.NewUserRepository(db)
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, jwtService)

	router := gin.Default()

	routes.SetupRoutes(router, authService, jwtService)

	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
