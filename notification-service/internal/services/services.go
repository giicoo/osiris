package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giicoo/osiris/notification-service/internal/config"
	"github.com/giicoo/osiris/notification-service/internal/entity"
	"github.com/giicoo/osiris/notification-service/internal/infrastructure"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Services struct {
	ctx            context.Context
	cfg            *config.Config
	sessionManager *ManagerSession
	rabbitMQ       infrastructure.AlertProducing
}

func NewServices(cfg *config.Config, rabbitMQ infrastructure.AlertProducing) *Services {
	return &Services{
		ctx:            context.Background(),
		cfg:            cfg,
		sessionManager: NewManagerSession(),
		rabbitMQ:       rabbitMQ,
	}
}

func (s *Services) Processing(msg []byte) error {
	var process entity.Process
	if err := json.Unmarshal(msg, &process); err != nil {
		return fmt.Errorf("unmarshal alert request: %w", err)
	}
	logrus.Info(process)

	session, err := s.sessionManager.GetSessionConn(2136636238)
	if err != nil {
		return fmt.Errorf("get session conn: %w", err)
	}
	if err := session.WriteJSON(process); err != nil {
		return fmt.Errorf("send json ws: %w", err)
	}
	return nil
}

func (s *Services) CreateSession(userID int, conn *websocket.Conn) {
	s.sessionManager.CreateSession(userID, conn)
}
