package service

import (
	"crypto/sha1"
	"filetranslation/pkg/models"
	"filetranslation/pkg/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	salt       = "dadadadada"
	tokenTTL   = 12 * time.Hour
	signingKey = "dadadda" // исправлено: singingKey -> signingKey
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"` // исправлено: UserId -> UserID
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}