package entities

import (
	"github.com/google/uuid"
	"time"
)

type Tweet struct {
	ID      uuid.UUID    `json:"id" validate:"required,uuid4"`
	UserID  uuid.UUID    `json:"userID" validate:"required,uuid4"`
	Tweet   string       `json:"tweet" validate:"required"`
	Created time.Time    `json:"created" validate:"required"` // time.Unix ?
	Likes   *[]uuid.UUID `json:"likes" validate:"required"`
}
