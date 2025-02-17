package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/giicoo/osiris/points-service/internal/config"
	"github.com/giicoo/osiris/points-service/internal/entity"
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
	tablesSQL, err := os.ReadFile("internal/repository/postgres/sql/create_tables.sql")
	if err != nil {
		return fmt.Errorf("open sql: %w", err)
	}
	if _, err := repo.db.Exec(string(tablesSQL)); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}
	return nil
}

func (repo *Repo) CloseConnection() error {
	return repo.db.Close()
}

func (repo *Repo) CreatePoint(point *entity.Point) (int, error) {
	var id int
	rows := repo.db.QueryRow(`INSERT INTO points (user_id, title, location, radius) VALUES ($1, $2, ST_GeomFromText($3, 4326), $4) RETURNING id;`, point.UserID, point.Title, point.Location, point.Radius)
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("db create scan point: %w", err)
	}
	return id, nil
}

func (repo *Repo) GetPoint(id int) (*entity.Point, error) {
	var point entity.Point
	if err := repo.db.QueryRow("SELECT id, user_id, title, ST_AsText(location), radius, created_at, updated_at FROM points WHERE id=$1", id).Scan(&point.ID, &point.UserID, &point.Title, &point.Location, &point.Radius, &point.CreatedAt, &point.UpdatedAt); err != nil {
		return nil, fmt.Errorf("db get point: %w", err)
	}

	return &point, nil
}

func (repo *Repo) GetPoints(user_id int) ([]*entity.Point, error) {
	rows, err := repo.db.Query("SELECT id, user_id, title, ST_AsText(location), radius, created_at, updated_at FROM points WHERE user_id=$1", user_id)
	if err != nil {
		return nil, fmt.Errorf("db get points: %w", err)
	}
	defer rows.Close()

	var points []*entity.Point
	for rows.Next() {
		var point entity.Point
		if err := rows.Scan(&point.ID, &point.UserID, &point.Title, &point.Location, &point.Radius, &point.CreatedAt, &point.UpdatedAt); err != nil {
			return nil, fmt.Errorf("db scan point: %w", err)
		}
		points = append(points, &point)
	}
	logrus.Info(points)
	return points, nil
}

func (repo *Repo) DeletePoint(id int) error {
	if _, err := repo.db.Query("DELETE FROM points WHERE id=$1", id); err != nil {
		return fmt.Errorf("db delete point: %w", err)
	}

	return nil
}

func (repo *Repo) UpdateTitle(id int, title string) error {
	if _, err := repo.db.Query("UPDATE points SET title=$2, updated_at=DEFAULT WHERE id=$1", id, title); err != nil {
		return fmt.Errorf("db update point: %w", err)
	}
	return nil
}

func (repo *Repo) UpdateLocation(id int, location string) error {
	if _, err := repo.db.Query("UPDATE points SET location=ST_GeomFromText($2, 4326), updated_at=DEFAULT  WHERE id=$1", id, location); err != nil {
		return fmt.Errorf("db update point: %w", err)
	}
	return nil
}

func (repo *Repo) UpdateRadius(id int, radius int) error {
	if _, err := repo.db.Query("UPDATE points SET radius=$2, updated_at=DEFAULT  WHERE id=$1", id, radius); err != nil {
		return fmt.Errorf("db update point: %w", err)
	}
	return nil
}
