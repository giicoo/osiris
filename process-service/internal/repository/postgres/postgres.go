package postgres

import (
	"database/sql"
	"fmt"

	"github.com/giicoo/osiris/process-service/internal/config"
	"github.com/giicoo/osiris/process-service/internal/entity"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	cfg *config.Config
	db  *sql.DB
}

func NewRepo(cfg *config.Config) *Repo {
	return &Repo{
		cfg: cfg,
	}
}

func (repo *Repo) Connection() error {
	connStr := fmt.Sprintf("postgres://%s:%s@db/db?sslmode=disable", repo.cfg.Database.User, repo.cfg.Database.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	repo.db = db

	if err := repo.db.Ping(); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	return nil
}

func (repo *Repo) CloseConnection() error {
	return repo.db.Close()
}

func (repo *Repo) GetNeedPoints(id int) ([]entity.Point, error) {
	logrus.Info(id)
	rows, err := repo.db.Query("SELECT p.* FROM points p JOIN alerts a ON ST_DWithin(a.location::geography, p.location::geography, a.radius + p.radius) WHERE a.id=$1;", id)
	if err != nil {
		return nil, fmt.Errorf("db get points: %w", err)
	}
	defer rows.Close()

	var points []entity.Point
	for rows.Next() {
		var point entity.Point
		if err := rows.Scan(&point.ID, &point.UserID, &point.Title, &point.Location, &point.Radius, &point.CreatedAt, &point.UpdatedAt); err != nil {
			return nil, fmt.Errorf("db scan point: %w", err)
		}
		points = append(points, point)
	}
	return points, nil
}
