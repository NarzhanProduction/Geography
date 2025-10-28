package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	request_models "github.com/NarzhanProduction/Geography/internal/api/models"
	"github.com/NarzhanProduction/Geography/internal/database/models"
	"github.com/NarzhanProduction/Geography/internal/database/repository"
	"github.com/NarzhanProduction/Geography/internal/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger logger.Logger
	repo   repository.UserRepository
	jwtKey []byte
	jwtTTL time.Duration
}

func NewUserService(logger logger.Logger, repo repository.UserRepository, jwtKey string, jwtTTL time.Duration) *UserService {
	return &UserService{
		logger: logger,
		repo:   repo,
		jwtKey: []byte(jwtKey),
		jwtTTL: jwtTTL,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *UserService) CreateUser(ctx context.Context, req *request_models.UserCreateRequest) error {
	hash, err := hashPassword(req.Password)
	if err != nil {
		return err
	}
	user := &models.User{Name: req.Name, Email: req.Email, PasswordHash: hash}
	return s.repo.Create(ctx, user)
}

func (s *UserService) LoginCheck(ctx context.Context, req request_models.UserLoginRequest) (*request_models.LoginResponse, error) {
	user, err := s.repo.GetByLogin(ctx, req.Name)
	if err != nil {
		s.logger.Error(ctx, fmt.Sprintf("error getting login: %v", err))
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Error(ctx, fmt.Sprintf("password of %s doesn't match: %s and %s", user.Name, user.PasswordHash, req.Password))
		return nil, errors.New("password doesn't match")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.jwtTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		s.logger.Error(ctx, fmt.Sprintf("failed to generate token: %v", err))
		return nil, errors.New("failed to generate token")
	}

	return &request_models.LoginResponse{Token: tokenString}, nil
}
