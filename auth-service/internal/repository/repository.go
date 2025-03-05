package repository

import "github.com/giicoo/osiris/auth-service/internal/entity"

type Repo interface {
	CreateUser(user *entity.User) (int, error)
	// GetUser(id int) (*entity.User, error)
}
