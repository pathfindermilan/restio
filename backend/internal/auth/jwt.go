package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
}

func NewJWTService(secret string) JWTService {
	return &jwtService{
		secretKey: secret,
	}
}

func (j *jwtService) GenerateToken(userID uint) (string, error) {
	claims := &jwtCustomClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtCustomClaim); !ok || !token.Valid || claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("invalid or expired token")
	}

	return token, nil
}
