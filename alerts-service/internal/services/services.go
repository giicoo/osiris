package services

import (
	"github.com/giicoo/osiris/alerts-service/internal/config"
	"github.com/giicoo/osiris/alerts-service/internal/repository"
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
