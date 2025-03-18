package entity

type Process struct {
	Alert  Alert   `json:"alert"`
	Points []Point `json:"points"`
}
