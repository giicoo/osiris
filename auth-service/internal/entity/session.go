package entity

type Session struct {
	ID     string `redis:"id"`
	UserID int    `redis:"user_id"`
}
