package postgres

import (
	"database/sql"
	"fmt"

	"github.com/giicoo/osiris/auth-service/internal/config"
	"github.com/giicoo/osiris/auth-service/internal/entity"

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

func (repo *Repo) CreateUser(user *entity.User) (int, error) {
	
	var id int
	rows := repo.db.QueryRow(`INSERT INTO users (id, login_yandex, first_name, last_name) VALUES ($1, $2, $3,$4) ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id, first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name RETURNING id;`, user.ID, user.LoginYandex, user.FirstName, user.LastName)
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("db create scan user: %w", err)
	}

	return id, nil
}

func (repo *Repo) GetUser(id int) (*entity.User, error) {
	var user entity.User
	if err := repo.db.QueryRow("SELECT * FROM users WHERE id=$1", id).Scan(&user.ID, &user.LoginYandex, &user.FirstName, &user.LastName); err != nil {
		return nil, fmt.Errorf("db get user: %w", err)
	}

	return &user, nil
}
