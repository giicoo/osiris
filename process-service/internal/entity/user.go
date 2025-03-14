package entity

type User struct {
	ID          int    `json:"id"`
	LoginYandex string `json:"login"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}
