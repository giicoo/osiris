package repository

import "github.com/giicoo/osiris/points-service/internal/entity"

type Repo interface {
	Connection() error
	CloseConnection() error
	CreatePoint(point *entity.Point) (int, error)
	GetPoint(id int) (*entity.Point, error)
	GetPoints(user_id int) ([]*entity.Point, error)
	DeletePoint(id int) error
	UpdateTitle(id int, title string) error
	UpdateLocation(id int, location string) error
	UpdateRadius(id int, radius int) error
}
