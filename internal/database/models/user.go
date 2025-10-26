package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	IsAdmin      bool
}
