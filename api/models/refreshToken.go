package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID          string    `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Fingerprint string    `json:"fingerprint"`
	CreatedAt   time.Time `json:"createdAt"`
}
