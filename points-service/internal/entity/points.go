package entity

import (
	"time"
)

type Point struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Location  string    `json:"location"`
	Radius    int       `json:"radius"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
