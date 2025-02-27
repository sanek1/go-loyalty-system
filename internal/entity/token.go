package entity

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID           uuid.UUID `pg:"type:uuid,pk"`
	UserID       uint
	CreationDate time.Time
	UsedAt       time.Time
}
