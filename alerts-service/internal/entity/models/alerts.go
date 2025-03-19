package models

type CreateAlert struct {
	Title       string `json:"title" binding:"required"  example:"Test Alert"`
	Description string `json:"description" binding:"required" example:"About Alert"`
	TypeID      int    `json:"type_id" binding:"required" example:"1"`
	Location    string `json:"location" binding:"required" example:"POINT(0 0)"`
	Radius      int    `json:"radius" binding:"required" example:"1"`
	Status      bool   `json:"status" binding:"required" example:"true"`
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
