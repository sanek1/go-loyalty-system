package entity

import (
	"time"
)

type OrderStatus string
type OrderStatusID int

type Order struct {
	ID           uint          `json:"ID"`
	UserID       uint          `json:"USER_ID"`
	StatusID     OrderStatusID `json:"StatusId"`
	CreationDate string        `json:"CreationDate"`
	Number       string        `json:"Number"`
	CreatedAt    time.Time     `json:"CreatedAt"`
	UploadedAt   time.Time     `json:"Uploaded"`
}

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

const (
	OrderStatusNewID        OrderStatusID = 1
	OrderStatusProcessingID OrderStatusID = 2
	OrderStatusInvalidID    OrderStatusID = 3
	OrderStatusProcessedID  OrderStatusID = 4
)

type OrderResponse struct {
	ID         uint      `json:"id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type OrderResponseDto struct {
	Number     int       `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}
