package entity

type Accrual struct {
	ID                  int64   `json:"id"`
	Accrual_statuses_id int64   `json:"accrual_statuses_id"`
	Accrual             float64 `json:"accrual"`
}
