package session

import (
	"context"
	"fmt"

	"github.com/giicoo/osiris/auth-service/internal/entity"
	"github.com/giicoo/osiris/auth-service/internal/repository"
	"github.com/giicoo/osiris/auth-service/internal/repository/redisRepo"
)

type SessionManager struct {
	ctx  context.Context
	repo repository.SessionRepo
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		ctx:  context.Background(),
		repo: redisRepo.NewSessionRepo(),
	}
}

func (sm *SessionManager) CreateSession(s *entity.Session) (*entity.Session, error) {
	id, err := sm.repo.CreateSession(sm.ctx, s)
	if err != nil {
		return nil, fmt.Errorf("session manager create: %w", err)
	}

	return id, nil
}

func (sm *SessionManager) GetSession(id string) (*entity.Session, error) {
	session, err := sm.repo.GetSession(sm.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("session manager get '%s': %w", id, err)
	}

	return session, nil
}

func (sm *SessionManager) DeleteSession(id string) error {
	session, err := sm.GetSession(id)
	if err != nil {
		return fmt.Errorf("session manager delete '%s': %w", id, err)
	}
	if err := sm.repo.DeleteSession(sm.ctx, id); err != nil {
		return fmt.Errorf("session manager delete '%s': %w", id, err)
	}
	if err := sm.repo.DeleteSessionFromUser(sm.ctx, id, session.UserID); err != nil {
		return fmt.Errorf("session manager delete '%s': %w", id, err)
	}
	return nil
}

func (sm *SessionManager) GetListSession(user_id int) ([]*entity.Session, error) {
	sessions, err := sm.repo.GetListSession(sm.ctx, user_id)
	if err != nil {
		return nil, fmt.Errorf("session manager get list '%d': %w", user_id, err)
	}
	return sessions, err
}

func (sm *SessionManager) DeleteListSession(user_id int) error {
	if err := sm.repo.DeleteListSession(sm.ctx, user_id); err != nil {
		return fmt.Errorf("session manager del list '%d': %w", user_id, err)
	}

	return nil
}
