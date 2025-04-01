package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/giicoo/osiris/auth-service/internal/config"
	"github.com/giicoo/osiris/auth-service/internal/entity"
	"github.com/giicoo/osiris/auth-service/internal/repository"
	"github.com/giicoo/osiris/auth-service/internal/services/session"
	"github.com/giicoo/osiris/auth-service/pkg/apiError"
	hashTools "github.com/giicoo/osiris/auth-service/pkg/hash"
)

type Services struct {
	cfg            *config.Config
	repo           repository.Repo
	sessionManager *session.SessionManager
}

func NewServices(cfg *config.Config, repo repository.Repo) *Services {
	return &Services{
		cfg:            cfg,
		repo:           repo,
		sessionManager: session.NewSessionManager(),
	}
}

func (s *Services) CreateUser(user *entity.User) (*entity.User, apiError.AErr) {
	userYet, err := s.repo.GetUserByEmail(user.Email)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return nil, apiError.New(fmt.Errorf("check exist user: %w", err), http.StatusInternalServerError)
	}
	if userYet != nil {
		return nil, apiError.New(fmt.Errorf("user already exist"), http.StatusBadRequest)
	}

	hashPassword, err := hashTools.HashPassword(user.Password)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("hash password: %w", err), http.StatusInternalServerError)
	}
	user.Password = hashPassword

	userID, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("create user: %w", err), http.StatusInternalServerError)
	}

	userDB, err := s.repo.GetUserById(userID)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("get user by %s email: %w", user.Email, err), http.StatusInternalServerError)
	}
	return userDB, nil
}

func (s *Services) CheckUser(user *entity.User) (string, apiError.AErr) {
	userDB, err := s.repo.GetUserByEmail(user.Email)
	if errors.Is(err, sql.ErrNoRows) && err != nil {
		return "", apiError.New(fmt.Errorf("user not exist"), http.StatusBadRequest)
	}
	if err != nil {
		return "", apiError.New(fmt.Errorf("user check %s: %w", user.Email, err), http.StatusInternalServerError)
	}

	if !hashTools.CheckPasswordHash(user.Password, userDB.Password) {
		return "", apiError.New(fmt.Errorf("wrong password"), http.StatusBadRequest)
	}
	session := &entity.Session{UserID: userDB.ID}
	sessionDB, err := s.sessionManager.CreateSession(session)
	if err != nil {
		return "", apiError.New(fmt.Errorf("create session: %w", err), http.StatusInternalServerError)
	}
	return sessionDB.ID, nil
}

func (s *Services) GetUserByID(id int) (*entity.User, apiError.AErr) {
	userDB, err := s.repo.GetUserById(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apiError.New(fmt.Errorf("user not exist"), http.StatusBadRequest)
	}
	if err != nil {
		return nil, apiError.New(fmt.Errorf("get user by %d id: %w", id, err), http.StatusInternalServerError)
	}
	return userDB, nil
}

func (s *Services) GetUserFromSession(sessionID string) (*entity.User, apiError.AErr) {
	session, err := s.sessionManager.GetSession(sessionID)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("get session: %w", err), http.StatusInternalServerError)
	}

	user, err := s.repo.GetUserById(session.UserID)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("get user: %w", err), http.StatusInternalServerError)
	}
	return user, nil
}

func (s *Services) DeleteSession(sessionID string) error {
	return s.sessionManager.DeleteSession(sessionID)
}
