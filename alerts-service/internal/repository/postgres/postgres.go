package postgres

import (
	"database/sql"
	"fmt"

	"github.com/giicoo/osiris/alerts-service/internal/config"
	"github.com/giicoo/osiris/alerts-service/internal/entity"
	_ "github.com/lib/pq"
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

func (repo *Repo) CreateAlert(alert *entity.Alert) (int, error) {
	var id int
	rows := repo.db.QueryRow(`INSERT INTO alerts ("user_id", "title", "description", "type_id", "location", "radius", "status")  VALUES ($1, $2, $3, $4, ST_GeomFromText($5, 4326), $6, $7) RETURNING id;`, alert.UserID, alert.Title, alert.Description, alert.TypeID, alert.Location, alert.Radius, alert.Status)
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("db create scan alert: %w", err)
	}
	return id, nil
}

func (repo *Repo) CreateType(typeModel *entity.Type) (int, error) {
	var id int
	rows := repo.db.QueryRow(`INSERT INTO types ("title")  VALUES ($1) RETURNING id;`, typeModel.Title)
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("db create scan type: %w", err)
	}
	return id, nil
}

func (repo *Repo) GetAlert(id int) (*entity.Alert, error) {
	var alert entity.Alert
	if err := repo.db.QueryRow(`SELECT id, user_id, title, description, type_id, ST_AsText(location), radius, status, created_at, updated_at FROM alerts WHERE id=$1`, id).Scan(&alert.ID, &alert.UserID, &alert.Title, &alert.Description, &alert.TypeID, &alert.Location, &alert.Radius, &alert.Status, &alert.CreatedAt, &alert.UpdatedAt); err != nil {
		return nil, fmt.Errorf("db get alert: %w", err)
	}

	return &alert, nil
}

func (repo *Repo) GetType(id int) (*entity.Type, error) {
	var typeModel entity.Type
	if err := repo.db.QueryRow("SELECT id, title FROM types WHERE id=$1", id).Scan(&typeModel.ID, &typeModel.Title); err != nil {
		return nil, fmt.Errorf("db get type: %w", err)
	}

	return &typeModel, nil
}

func (repo *Repo) GetAlerts() ([]*entity.Alert, error) {
	rows, err := repo.db.Query("SELECT id, user_id, title, description, type_id, ST_AsText(location), radius, status, created_at, updated_at FROM alerts")
	if err != nil {
		return nil, fmt.Errorf("db get alerts: %w", err)
	}
	defer rows.Close()

	var alerts []*entity.Alert
	for rows.Next() {
		var alert entity.Alert
		if err := rows.Scan(&alert.ID, &alert.UserID, &alert.Title, &alert.Description, &alert.TypeID, &alert.Location, &alert.Radius, &alert.Status, &alert.CreatedAt, &alert.UpdatedAt); err != nil {
			return nil, fmt.Errorf("db scan alerts: %w", err)
		}
		alerts = append(alerts, &alert)
	}
	return alerts, nil
}

func (repo *Repo) GetTypes() ([]*entity.Type, error) {
	rows, err := repo.db.Query("SELECT id, title FROM types")
	if err != nil {
		return nil, fmt.Errorf("db get types: %w", err)
	}
	defer rows.Close()

	var types []*entity.Type
	for rows.Next() {
		var typeModel entity.Type
		if err := rows.Scan(&typeModel.ID, &typeModel.Title); err != nil {
			return nil, fmt.Errorf("db scan types: %w", err)
		}
		types = append(types, &typeModel)
	}
	return types, nil
}

func (repo *Repo) DeleteType(id int) error {
	if _, err := repo.db.Query("DELETE FROM types WHERE id=$1", id); err != nil {
		return fmt.Errorf("db delete type: %w", err)
	}
	return nil
}
