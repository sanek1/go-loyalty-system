package entity

import "errors"

var (
	ErrInvalidOrder = errors.New("invalid order number")
	ErrOrderExists  = errors.New("order already exists")

	ErrInvalidUser     = errors.New("invalid user")
	ErrUserExists      = errors.New("user already exists")
	FailedToCheckOrder = errors.New("failed to check order existence")

	UserDoesNotExist      = errors.New("user does not exist")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidOrderNumber = errors.New("invalid order number")

	ErrOrderExistsThisUser  = errors.New("order already uploaded by this user")
	ErrOrderExistsOtherUser = errors.New("order already uploaded by another user")
)
