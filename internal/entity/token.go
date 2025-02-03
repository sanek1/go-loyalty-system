package entity

import (
	"time"
)

type Token struct {
	ID           uint
	UserID       uint
	CreationDate time.Time
	UsedAt       time.Time
}
