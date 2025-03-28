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
	rows := repo.db.QueryRow(`INSERT INTO users (email, password)  VALUES ($1, $2) RETURNING id;`, user.Email, user.Password)
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("db create scan user: %w", err)
	}
	return id, nil
}

func (repo *Repo) GetUserById(id int) (*entity.User, error) {
	var user entity.User
	if err := repo.db.QueryRow("SELECT id, email, password FROM users WHERE id=$1", id).Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return nil, fmt.Errorf("db get user: %w", err)
	}

	return &user, nil
}

func (repo *Repo) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := repo.db.QueryRow("SELECT id, email, password FROM users WHERE email=$1", email).Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return nil, fmt.Errorf("db get user: %w", err)
	}

	return &user, nil
}
