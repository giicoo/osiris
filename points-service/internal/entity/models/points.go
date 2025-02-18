package models

type CreatePoint struct {
	UserID   int    `extensions:"x-order=0" json:"user_id" binding:"required"`
	Title    string `extensions:"x-order=1" json:"title" binding:"required"`
	Location string `extensions:"x-order=2" json:"location" binding:"required"`
	Radius   int    `extensions:"x-order=3" json:"radius" binding:"required"`
}

type UpdateTitlePoint struct {
	ID    int    `extensions:"x-order=0" json:"id" binding:"required"`
	Title string `extensions:"x-order=1" json:"title" binding:"required"`
}

type UpdateLocationPoint struct {
	ID       int    `extensions:"x-order=0" json:"id" binding:"required"`
	Location string `extensions:"x-order=1" json:"location" binding:"required"`
}

type UpdateRadiusPoint struct {
	ID     int `extensions:"x-order=0" json:"id" binding:"required"`
	Radius int `extensions:"x-order=1" json:"radius" binding:"required"`
}

type DeletePoint struct {
	ID int `extensions:"x-order=0" json:"id" binding:"required"`
}
