package services

import (
	"backend/internal/config"
	"backend/internal/workers"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-co-op/gocron"

	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *models.User) error
	Login(identifier string, password string) (string, error)
	VerifyEmail(username string, verificationCode string) error
	DeleteUser(userID uint) error
	UpdateUser(userID uint, updatedData *models.User) error
	GetProfile(userID uint) (*models.User, error)
	Logout() error
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

	user.VerificationCode = ""
	user.IsVerified = true
	return s.userRepo.UpdateUser(user)
}

func (s *authService) DeleteUser(userID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.DeletedAt != nil {
		return errors.New("user already deleted")
	}

	if err := s.userRepo.DeleteUser(user); err != nil {
		return err
	}

	// JWT invalidation logic goes here if needed

	return nil
}

func (s *authService) UpdateUser(userID uint, updatedData *models.User) error {
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return err
	}

	if err := s.userRepo.UpdateUser(updatedData); err != nil {
		return err
	}

	return nil
}

func (s *authService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Logout() error {
	// JWT invalidation logic goes here if needed
	return nil
}
