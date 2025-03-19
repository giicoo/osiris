package entity

import "time"

type Error struct {
	Error string `json:"error"`
}

type Response struct {
	Response string `json:"message"`
}
type Alert struct {
	ID          int       `extensions:"x-order=0" json:"id"`
	UserID      int       `extensions:"x-order=1" json:"user_id"`
	Title       string    `extensions:"x-order=2" json:"title"`
	Description string    `extensions:"x-order=3" json:"description"`
	TypeID      int       `extensions:"x-order=4" json:"type_id"`
	Location    string    `extensions:"x-order=5" json:"location"`
	Radius      int       `extensions:"x-order=6" json:"radius"`
	Status      bool      `extensions:"x-order=7" json:"status"`
	CreatedAt   time.Time `extensions:"x-order=8" json:"created_at"`
	UpdatedAt   time.Time `extensions:"x-order=9" json:"updated_at"`
}

type Type struct {
	ID    int    `extensions:"x-order=0" json:"id"`
	Title string `extensions:"x-order=1" json:"title"`
}
