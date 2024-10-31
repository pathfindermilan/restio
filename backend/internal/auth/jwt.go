package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	InvalidateToken(token string)
	IsTokenInvalidated(token string) bool
}

type jwtCustomClaim struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey    string
	blacklist    map[string]struct{}
	blacklistMux sync.RWMutex
}

func NewJWTService(secret string) JWTService {
	return &jwtService{
		secretKey: secret,
		blacklist: make(map[string]struct{}),
	}
}

func (j *jwtService) GenerateToken(userID uint) (string, error) {
	claims := &jwtCustomClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Id:        generateUniqueTokenID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func generateUniqueTokenID() string {
	return uuid.New().String()
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtCustomClaim)
	if !ok || !token.Valid || claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("invalid or expired token")
	}

	if j.IsTokenInvalidated(claims.Id) {
		return nil, errors.New("token is invalidated")
	}

	return token, nil
}

func (j *jwtService) InvalidateToken(tokenString string) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*jwtCustomClaim); ok && token.Valid {
		j.blacklistMux.Lock()
		defer j.blacklistMux.Unlock()
		j.blacklist[claims.Id] = struct{}{}
	}
}

func (j *jwtService) IsTokenInvalidated(tokenID string) bool {
	j.blacklistMux.RLock()
	defer j.blacklistMux.RUnlock()
	_, exists := j.blacklist[tokenID]
	return exists
}
