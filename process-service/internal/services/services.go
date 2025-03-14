package services

import (
	"encoding/json"
	"fmt"

	"github.com/giicoo/osiris/process-service/internal/config"
	"github.com/giicoo/osiris/process-service/internal/entity"
	"github.com/giicoo/osiris/process-service/internal/infrastructure"
	"github.com/giicoo/osiris/process-service/internal/repository"
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

func (s *Services) Processing(msg []byte) error {
	var alert entity.Alert
	if err := json.Unmarshal(msg, &alert); err != nil {
		return fmt.Errorf("unmarshal alert request: %w", err)
	}
	points, err := s.repo.GetNeedPoints(alert.ID)
	if err != nil {
		return fmt.Errorf("get need points: %w", err)
	}
	body, err := json.Marshal(points)
	if err != nil {
		return fmt.Errorf("marshal points: %w", err)
	}
	s.rabbitMQ.PublicMessage(body)
	return nil
}
