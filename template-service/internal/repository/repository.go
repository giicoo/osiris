package repository

import "github.com/osiris/template-service/internal/config"

type Repo interface {
}

type RepoTemp struct {
	cfg *config.Config
}

func NewRepoTemp(cfg *config.Config) Repo {
	return &RepoTemp{cfg: cfg}
}
