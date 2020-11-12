package entity

import "time"

type Currency struct {
	Title     string    `json:"title"`
	Value     float64   `json:"value,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Base      string    `json:"base,omitempty"`
}
