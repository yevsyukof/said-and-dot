package entities

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id" validate:"required,uuid4"`
	Username  string    `json:"username" validate:"required"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
}
