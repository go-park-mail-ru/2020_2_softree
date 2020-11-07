package entity

import "time"

type Currency struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Value     float64   `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
	Base      string
}
