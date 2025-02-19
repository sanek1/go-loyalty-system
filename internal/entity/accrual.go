package entity

type Accrual struct {
	ID                int64   `json:"id"`
	AccrualStatusesID int64   `json:"accrual_statuses_id"`
	Accrual           float64 `json:"accrual"`
}

type AccrualOrder struct {
	Order string    `json:"order"`
	Goods []Product `json:"goods"`
}

type Product struct {
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
