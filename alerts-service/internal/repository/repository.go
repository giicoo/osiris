package repository

import "github.com/giicoo/osiris/alerts-service/internal/entity"

type Repo interface {
	Connection() error
	CloseConnection() error
	CreateAlert(alert *entity.Alert) (int, error)
	UpdateStatusAlert(id int, status bool) error
	CreateType(typeModel *entity.Type) (int, error)
	GetAlert(id int) (*entity.Alert, error)
	GetType(id int) (*entity.Type, error)
	GetAlerts() ([]*entity.Alert, error)
	GetTypes() ([]*entity.Type, error)
	DeleteType(id int) error
}
