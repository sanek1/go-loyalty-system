package entity

import "errors"

var (
	ErrInvalidOrder = errors.New("invalid order number") //422
	ErrOrderExists  = errors.New("order already exists")

	ErrInvalidUser        = errors.New("invalid user")
	ErrUserExists         = errors.New("user already exists")
	ErrFailedToCheckOrder = errors.New("failed to check order existence")

	ErrUserDoesNotExist   = errors.New("user does not exist")  //401
	ErrInsufficientFunds  = errors.New("insufficient funds")   //402
	ErrInvalidOrderNumber = errors.New("invalid order number") //422

	ErrOrderExistsThisUser  = errors.New("order already uploaded by this user")    //200
	ErrOrderExistsOtherUser = errors.New("order already uploaded by another user") //200
)
