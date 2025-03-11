package services

import (
	"fmt"
	"net/http"

	"github.com/giicoo/osiris/points-service/internal/config"
	"github.com/giicoo/osiris/points-service/internal/entity"
	"github.com/giicoo/osiris/points-service/internal/repository"
	"github.com/giicoo/osiris/points-service/pkg/apiError"
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

func (s *Services) CreatePoint(point *entity.Point) (*entity.Point, apiError.AErr) {
	id, err := s.repo.CreatePoint(point)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service create point: %w", err), http.StatusInternalServerError)
	}

	pointDB, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service create get point: %w", err), http.StatusInternalServerError)
	}
	return pointDB, nil
}

func (s *Services) GetPoint(id int, userID int) (*entity.Point, apiError.AErr) {
	// check if a point belongs to a user
	point, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get point: %w", err), http.StatusInternalServerError)
	}
	if point.UserID != userID {
		return nil, apiError.New(fmt.Errorf("point not by user with %d id", id), http.StatusBadRequest)
	}
	return point, nil
}

func (s *Services) GetPoints(userID int) ([]*entity.Point, apiError.AErr) {
	points, err := s.repo.GetPoints(userID)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get point: %w", err), http.StatusInternalServerError)
	}
	return points, nil
}

func (s *Services) DeletePoint(id int, userID int) apiError.AErr {
	// check if a point belongs to a user
	point, err := s.repo.GetPoint(id)
	if err != nil {
		return apiError.New(fmt.Errorf("service get point: %w", err), http.StatusInternalServerError)
	}
	if point.UserID != userID {
		return apiError.New(fmt.Errorf("point not by user with %d id", id), http.StatusBadRequest)
	}

	if err := s.repo.DeletePoint(id); err != nil {
		return apiError.New(fmt.Errorf("service delete point: %w", err), http.StatusInternalServerError)
	}
	return nil
}

func (s *Services) UpdateTitle(id int, title string, userID int) (*entity.Point, apiError.AErr) {
	// check if a point belongs to a user
	point, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get point: %w", err), http.StatusInternalServerError)
	}
	if point.UserID != userID {
		return nil, apiError.New(fmt.Errorf("point not by user with %d id", id), http.StatusBadRequest)
	}

	if err := s.repo.UpdateTitle(id, title); err != nil {
		return nil, apiError.New(fmt.Errorf("service update point: %w", err), http.StatusInternalServerError)
	}

	pointDB, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service update get point: %w", err), http.StatusInternalServerError)
	}

	return pointDB, nil
}

func (s *Services) UpdateLocation(id int, location string, userID int) (*entity.Point, apiError.AErr) {
	// check if a point belongs to a user
	point, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get point: %w", err), http.StatusInternalServerError)
	}
	if point.UserID != userID {
		return nil, apiError.New(fmt.Errorf("point not by user with %d id", id), http.StatusBadRequest)
	}

	if err := s.repo.UpdateLocation(id, location); err != nil {
		return nil, apiError.New(fmt.Errorf("service update point: %w", err), http.StatusInternalServerError)
	}

	pointDB, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get point: %w", err), http.StatusInternalServerError)
	}

	return pointDB, nil
}

func (s *Services) UpdateRadius(id int, radius int, userID int) (*entity.Point, apiError.AErr) {
	// check if a point belongs to a user
	point, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get point: %w", err), http.StatusInternalServerError)
	}
	if point.UserID != userID {
		return nil, apiError.New(fmt.Errorf("point not by user with %d id", id), http.StatusBadRequest)
	}
	if err := s.repo.UpdateRadius(id, radius); err != nil {
		return nil, apiError.New(fmt.Errorf("service update point: %w", err), http.StatusInternalServerError)
	}

	pointDB, err := s.repo.GetPoint(id)
	if err != nil {
		return nil, apiError.New(fmt.Errorf("service get point: %w", err), http.StatusInternalServerError)
	}

	return pointDB, nil
}
