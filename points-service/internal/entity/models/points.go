package models

type CreatePoint struct {
	UserID   int    `json:"user_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Location string `json:"location" binding:"required"`
	Radius   int    `json:"radius" binding:"required"`
}

type UpdateTitlePoint struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type UpdateLocationPoint struct {
	ID       int    `json:"id"`
	Location string `json:"location"`
}

type UpdateRadiusPoint struct {
	ID     int `json:"id"`
	Radius int `json:"radius"`
}

type DeletePoint struct {
	ID int `json:"id"`
}
