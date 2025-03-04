package repository

import "github.com/giicoo/osiris/auth-service/internal/config"

type Repo interface {
}

type RepoTemp struct {
	cfg *config.Config
}

func NewRepoTemp(cfg *config.Config) Repo {
	return &RepoTemp{cfg: cfg}
}
