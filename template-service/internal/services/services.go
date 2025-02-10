package services

import (
	"github.com/osiris/template-service/internal/config"
	"github.com/osiris/template-service/internal/repository"
)

type Services struct {
	cfg  *config.Config
	repo repository.Repo
}

func NewServices(cfg *config.Config, repo repository.Repo) *Services {
	return &Services{
		cfg:  cfg,
		repo: repo,
	}
}
