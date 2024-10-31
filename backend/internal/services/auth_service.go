package services

import (
	"backend/internal/config"
	"backend/internal/workers"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"math/rand"
	"strings"
	"time"

	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *models.User) error
	Login(identifier string, password string) (string, error)
	VerifyEmail(username string, verificationCode string) error
}

type authService struct {
	userRepo   repositories.UserRepository
	jwtService auth.JWTService
	scheduler  *gocron.Scheduler
}

func NewAuthService(userRepo repositories.UserRepository, jwtService auth.JWTService) AuthService {
	s := gocron.NewScheduler(time.UTC)
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
		scheduler:  s,
	}
}

func (s *authService) Register(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	code := rand.Intn(999999-100000) + 100000
	user.VerificationCode = fmt.Sprintf("%d", code)

	if err := s.userRepo.CreateUser(user); err != nil {
		return err
	}

	cfg := config.LoadConfig()

	s.scheduler.Every(1).Second().LimitRunsTo(1).Do(func() {
		err := workers.SendVerificationEmail(user, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.FrontendURL)
		if err != nil {
			fmt.Printf("Failed to send email: %s", err)
		}
	})
	s.scheduler.StartAsync()

	return nil
}

func (s *authService) Login(identifier, password string) (string, error) {
	var user *models.User
	var err error

	if strings.Contains(identifier, "@") {
		user, err = s.userRepo.GetUserByEmail(identifier)
	} else {
		user, err = s.userRepo.GetUserByUsername(identifier)
	}
	if err != nil {
		return "", err
	}

	if !user.IsVerified {
		return "", errors.New("email not verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) VerifyEmail(username, verificationCode string) error {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}

	if user.VerificationCode != verificationCode {
		return errors.New("invalid verification code")
	}

	user.VerificationCode = ""         // Clear the verification code after successful verification.
	user.IsVerified = true             // Set IsVerified to true after successful verification.
	return s.userRepo.UpdateUser(user) // Update user in repository.
}
