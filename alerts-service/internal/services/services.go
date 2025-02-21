package services

import (
	"encoding/json"
	"fmt"

	"github.com/giicoo/osiris/alerts-service/internal/config"
	"github.com/giicoo/osiris/alerts-service/internal/entity"
	"github.com/giicoo/osiris/alerts-service/internal/infrastructure/rabbitmq"
	"github.com/giicoo/osiris/alerts-service/internal/repository"
)

type Services struct {
	cfg      *config.Config
	repo     repository.Repo
	rabbitMQ *rabbitmq.AlertProducing
}

func NewServices(cfg *config.Config, repo repository.Repo, rabbitMQ *rabbitmq.AlertProducing) *Services {
	return &Services{
		cfg:      cfg,
		repo:     repo,
		rabbitMQ: rabbitMQ,
	}
}

func (s *Services) CreateAlert(alert *entity.Alert) (*entity.Alert, error) {
	id, err := s.repo.CreateAlert(alert)
	if err != nil {
		return nil, fmt.Errorf("service create alert: %w", err)
	}

	alertDB, err := s.repo.GetAlert(id)
	if err != nil {
		return nil, fmt.Errorf("service create get alert: %w", err)
	}
	body, err := json.Marshal(alert)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}

	if err := s.rabbitMQ.PublicMessage(body); err != nil {
		return nil, fmt.Errorf("publish message: %w", err)
	}
	return alertDB, nil
}

func (s *Services) CreateType(typeModel *entity.Type) (*entity.Type, error) {
	id, err := s.repo.CreateType(typeModel)
	if err != nil {
		return nil, fmt.Errorf("service create type: %w", err)
	}

	typeDB, err := s.repo.GetType(id)
	if err != nil {
		return nil, fmt.Errorf("service create get type: %w", err)
	}
	return typeDB, nil
}

func (s *Services) GetAlert(id int) (*entity.Alert, error) {
	alert, err := s.repo.GetAlert(id)
	if err != nil {
		return nil, fmt.Errorf("service get alert: %w", err)
	}
	return alert, nil
}

func (s *Services) GetType(id int) (*entity.Type, error) {
	typeModel, err := s.repo.GetType(id)
	if err != nil {
		return nil, fmt.Errorf("service get type: %w", err)
	}
	return typeModel, nil
}

func (s *Services) GetAlerts() ([]*entity.Alert, error) {
	alerts, err := s.repo.GetAlerts()
	if err != nil {
		return nil, fmt.Errorf("service get alerts: %w", err)
	}
	return alerts, nil
}

func (s *Services) GetTypes() ([]*entity.Type, error) {
	types, err := s.repo.GetTypes()
	if err != nil {
		return nil, fmt.Errorf("service get types: %w", err)
	}
	return types, nil
}

func (s *Services) DeleteType(id int) error {
	if err := s.repo.DeleteType(id); err != nil {
		return fmt.Errorf("service delete type: %w", err)
	}
	return nil
}
