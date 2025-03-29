package entity

type Process struct {
	Alert  Alert   `json:"alert"`
	Points []Point `json:"points"`
}

type ProcessUnic struct {
	Alert Alert `json:"alert"`
	Point Point `json:"point"`
}
