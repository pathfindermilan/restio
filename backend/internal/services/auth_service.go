package services

import (
	"errors"

	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *models.User) error
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo   repositories.UserRepository
	jwtService auth.JWTService
}

func NewAuthService(userRepo repositories.UserRepository, jwtService auth.JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Register(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.userRepo.CreateUser(user)
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
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
