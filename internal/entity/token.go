package entity

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID           uuid.UUID `pg:"type:uuid,pk"`
	UserID       uint
	CreationDate time.Time
	UsedAt       time.Time
}
