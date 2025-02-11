package entity

import "errors"

var (
	ErrInvalidOrder = errors.New("invalid order number")
	ErrOrderExists  = errors.New("order already exists")
)
