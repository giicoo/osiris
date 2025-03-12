package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/giicoo/osiris/alerts-service/internal/config"
	"github.com/giicoo/osiris/alerts-service/internal/entity"
	"github.com/giicoo/osiris/alerts-service/internal/infrastructure"
	"github.com/giicoo/osiris/alerts-service/internal/repository"
	"github.com/giicoo/osiris/alerts-service/pkg/apiError"
)

type Services struct {
	cfg      *config.Config
	repo     repository.Repo
	rabbitMQ infrastructure.AlertProducing
}

func NewServices(cfg *config.Config, repo repository.Repo, rabbitMQ infrastructure.AlertProducing) *Services {
	return &Services{
		cfg:      cfg,
		repo:     repo,
		rabbitMQ: rabbitMQ,
	}
}

func (s *Services) CreateAlert(alert *entity.Alert) (*entity.Alert, apiError.AErr) {
	id, err := s.repo.CreateAlert(alert)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service create alert: %w", err), http.StatusInternalServerError)
	}

	alertDB, err := s.repo.GetAlert(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service create get alert: %w", err), http.StatusInternalServerError)
	}
	body, err := json.Marshal(alert)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("failed to serialize message: %w", err), http.StatusInternalServerError)
	}

	if err := s.rabbitMQ.PublicMessage(body); err != nil {
		return nil, apiError.New(fmt.Errorf("publish message: %w", err), http.StatusInternalServerError)
	}
	return alertDB, nil
}

func (s *Services) CreateType(typeModel *entity.Type) (*entity.Type, apiError.AErr) {
	id, err := s.repo.CreateType(typeModel)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service create type: %w", err), http.StatusInternalServerError)
	}

	typeDB, err := s.repo.GetType(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service create type: %w", err), http.StatusInternalServerError)
	}
	return typeDB, nil
}

func (s *Services) GetAlert(id int) (*entity.Alert, apiError.AErr) {
	alert, err := s.repo.GetAlert(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get alert: %w", err), http.StatusInternalServerError)
	}
	return alert, nil
}

func (s *Services) GetType(id int) (*entity.Type, apiError.AErr) {
	typeModel, err := s.repo.GetType(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get type: %w", err), http.StatusInternalServerError)
	}
	return typeModel, nil
}

func (s *Services) GetAlerts() ([]*entity.Alert, apiError.AErr) {
	alerts, err := s.repo.GetAlerts()
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get alerts: %w", err), http.StatusInternalServerError)
	}
	return alerts, nil
}

func (s *Services) GetTypes() ([]*entity.Type, apiError.AErr) {
	types, err := s.repo.GetTypes()
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get types: %w", err), http.StatusInternalServerError)
	}
	return types, nil
}

func (s *Services) DeleteType(id int) apiError.AErr {
	if err := s.repo.DeleteType(id); err != nil {
		return apiError.New(fmt.Errorf("service delete type: %w", err), http.StatusInternalServerError)
	}
	return nil
}
