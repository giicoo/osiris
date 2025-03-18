package models

type CreateAlert struct {
	Title       string `extensions:"x-order=2" json:"title" binding:"required"`
	Description string `extensions:"x-order=3" json:"description" binding:"required"`
	TypeID      int    `extensions:"x-order=4" json:"type_id" binding:"required"`
	Location    string `extensions:"x-order=5" json:"location" binding:"required"`
	Radius      int    `extensions:"x-order=6" json:"radius" binding:"required"`
	Status      bool   `extensions:"x-order=7" json:"status" binding:"required"`
}

type DeleteAlert struct {
	ID int `extensions:"x-order=1" json:"id" binding:"required"`
}

type CreateType struct {
	Title string `extensions:"x-order=2" json:"title" binding:"required"`
}

type DeleteType struct {
	ID int `extensions:"x-order=1" json:"id" binding:"required"`
}

type StopAlert struct {
	ID int `extensions:"x-order=1" json:"id" binding:"required"`
}
