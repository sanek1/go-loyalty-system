package entity

type Balance struct {
	Current   float64 `json:"current_balance"`
	Withdrawn float64 `json:"withdrawn"`
}
