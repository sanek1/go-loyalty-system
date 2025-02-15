package entity

import "time"

type WithdrawalRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum" `
}

type Withdrawal struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	OrderNumber string    `json:"order"`
	Amount      float32   `json:"sum"`
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
