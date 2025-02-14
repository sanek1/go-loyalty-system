package entity

import (
	"time"

	"github.com/Azure/go-autorest/autorest/date"
)

type OrderStatus string
type OrderStatusID int

type Order struct {
	ID           uint      `json:"ID"`
	UserID       int       `json:"USER_ID"`
	StatusID     string    `json:"Status"`
	CreationDate string    `json:"CreationDate"`
	Number       string    `json:"Number"`
	UploadedAt   date.Time `json:"Uploaded"`
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
	Number     int       `json:"number"`
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
