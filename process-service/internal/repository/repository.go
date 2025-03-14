package repository

import "github.com/giicoo/osiris/process-service/internal/entity"

type Repo interface {
	Connection() error
	CloseConnection() error
	GetNeedPoints(id int) ([]entity.Point, error)
}
