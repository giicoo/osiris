package services

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ManagerSession struct {
	mu       sync.Mutex
	sessions map[int]*websocket.Conn
}

func NewManagerSession() *ManagerSession {
	return &ManagerSession{
		sessions: make(map[int]*websocket.Conn),
	}
}

func (m *ManagerSession) CreateSession(userID int, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[userID] = conn
	logrus.Info(userID, m.sessions)

}

func (m *ManagerSession) GetSessionConn(userID int) (*websocket.Conn, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// сессия сохраняется для юзера который создает alert
	logrus.Info(userID, m.sessions)
	conn, ok := m.sessions[userID]
	if !ok {
		return nil, fmt.Errorf("dont have session with user")
	}
	return conn, nil
}
