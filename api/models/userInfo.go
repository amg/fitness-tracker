package models

import (
	"github.com/google/uuid"
)

type UserInfo struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	PictureUrl string    `json:"pictureUrl"`
}
