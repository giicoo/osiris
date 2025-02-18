package services

import (
	"fmt"

	"github.com/giicoo/osiris/points-service/internal/config"
	"github.com/giicoo/osiris/points-service/internal/entity"
	"github.com/giicoo/osiris/points-service/internal/repository"
)

type Services struct {
	cfg  *config.Config
	repo repository.Repo
}

func NewServices(cfg *config.Config, repo repository.Repo) *Services {
	return &Services{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *Services) CreatePoint(point *entity.Point) (*entity.Point, error) {
	id, err := s.repo.CreatePoint(point)
	if err != nil {
		return nil, fmt.Errorf("service create point: %w", err)
	}

	pointDB, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, fmt.Errorf("service create get point: %w", err)
	}
	return pointDB, nil
}

func (s *Services) GetPoint(id int) (*entity.Point, error) {
	point, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, fmt.Errorf("service get point: %w", err)
	}
	return point, nil
}

func (s *Services) GetPoints(user_id int) ([]*entity.Point, error) {
	points, err := s.repo.GetPoints(user_id)
	if err != nil {
		return nil, fmt.Errorf("service get point: %w", err)
	}
	return points, nil
}

func (s *Services) DeletePoint(id int) error {
	if err := s.repo.DeletePoint(id); err != nil {
		return fmt.Errorf("service delete point: %w", err)
	}
	return nil
}

func (s *Services) UpdateTitle(id int, title string) (*entity.Point, error) {
	if err := s.repo.UpdateTitle(id, title); err != nil {
		return nil, fmt.Errorf("service update point: %w", err)
	}

	point, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, fmt.Errorf("service update get point: %w", err)
	}

	return point, nil
}

func (s *Services) UpdateLocation(id int, location string) (*entity.Point, error) {
	if err := s.repo.UpdateLocation(id, location); err != nil {
		return nil, fmt.Errorf("service update point: %w", err)
	}

	point, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, fmt.Errorf("service update get point: %w", err)
	}

	return point, nil
}

func (s *Services) UpdateRadius(id int, radius int) (*entity.Point, error) {
	if err := s.repo.UpdateRadius(id, radius); err != nil {
		return nil, fmt.Errorf("service update point: %w", err)
	}

	point, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, fmt.Errorf("service update get point: %w", err)
	}

	return point, nil
}
