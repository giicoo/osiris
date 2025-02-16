package entity

import (
	"time"
)

type Point struct {
	ID        int       `extensions:"x-order=0" json:"id"`
	UserID    int       `extensions:"x-order=1" json:"user_id"`
	Title     string    `extensions:"x-order=2" json:"title"`
	Location  string    `extensions:"x-order=3" json:"location"`
	Radius    int       `extensions:"x-order=4" json:"radius"`
	CreatedAt time.Time `extensions:"x-order=5" json:"created_at"`
	UpdatedAt time.Time `extensions:"x-order=6" json:"updated_at"`
}
