package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port         string
	JWTSecret    string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FrontendURL  string
}

func LoadConfig() Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	config := Config{
		Port:         viper.GetString("PORT"),
		JWTSecret:    viper.GetString("JWT_SECRET"),
		DBHost:       viper.GetString("DB_HOST"),
		DBPort:       viper.GetString("DB_PORT"),
		DBUser:       viper.GetString("DB_USER"),
		DBPassword:   viper.GetString("DB_PASSWORD"),
		DBName:       viper.GetString("DB_NAME"),
		SMTPHost:     viper.GetString("SMTP_HOST"),
		SMTPPort:     viper.GetString("SMTP_PORT"),
		SMTPUsername: viper.GetString("SMTP_USERNAME"),
		SMTPPassword: viper.GetString("SMTP_PASSWORD"),
		FrontendURL:  viper.GetString("FRONTEND_URL"),
	}

	return config
}
