package services

import (
	"fmt"

	"github.com/giicoo/osiris/auth-service/internal/config"
	"github.com/giicoo/osiris/auth-service/internal/entity"
	"github.com/giicoo/osiris/auth-service/internal/infrastructure"
	"github.com/giicoo/osiris/auth-service/internal/repository"
)

type Services struct {
	cfg  *config.Config
	repo repository.Repo

	ya *infrastructure.YandexAPI
}

func NewServices(cfg *config.Config, repo repository.Repo, ya *infrastructure.YandexAPI) *Services {
	return &Services{
		cfg:  cfg,
		repo: repo,
		ya:   ya,
	}
}

func (s *Services) CreateUser(code string) (string, error) {
	accessToken, err := s.ya.GetAccessToken(code)
	if err != nil {
		return "", fmt.Errorf("create user service: %w", err)
	}
	user, err := s.ya.GetUserInfo(accessToken)
	if err != nil {
		return "", fmt.Errorf("create user service: %w", err)
	}
	_, err = s.repo.CreateUser(user)
	if err != nil {
		return "", fmt.Errorf("create user db: %w", err)
	}

	return accessToken, nil
}

func (s *Services) GetUser(accessToken string) (*entity.User, error) {
	user, err := s.ya.GetUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("create user service: %w", err)
	}
	return user, nil
}
