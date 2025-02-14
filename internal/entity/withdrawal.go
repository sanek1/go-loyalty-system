package entity

import "time"

type WithdrawalRequest struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required,gt=0"`
}

type Withdrawal struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	OrderNumber string    `json:"order"`
	Amount      float64   `json:"sum"`
	CreatedAt   time.Time `json:"created_at"`
}
