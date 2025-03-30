package repository

import (
	"context"

	"github.com/giicoo/osiris/auth-service/internal/entity"
)

//go:generate mockgen -source=
type Repo interface {
	CreateUser(*entity.User) (int, error)
	GetUserById(id int) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
}

type SessionRepo interface {
	CreateSession(ctx context.Context, s *entity.Session) (*entity.Session, error)
	GetSession(ctx context.Context, id string) (*entity.Session, error)
	DeleteSession(ctx context.Context, id string) error
	DeleteSessionFromUser(ctx context.Context, id string, user_id int) error
	GetListSession(ctx context.Context, user_id int) ([]*entity.Session, error)
	DeleteListSession(ctx context.Context, user_id int) error
}
