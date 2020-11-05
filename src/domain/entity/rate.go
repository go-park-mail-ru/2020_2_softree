package entity

type Rate struct {
	ID       uint64 `json:"id"`
	Base     string `json:"base"`
	Currency string `json:"currency"`
	Value    float64 `json:"value"`
}
