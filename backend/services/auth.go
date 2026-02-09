package services

import (
	"context"
	"errors"
	"os"
	"time"

	"todo-app/dto"
	"todo-app/ent"
	"todo-app/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (string, error)
}

type AuthService struct {
	repo repositories.IUserRepository
}

func NewAuthService(repo repositories.IUserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // default for development/test if env not set
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
