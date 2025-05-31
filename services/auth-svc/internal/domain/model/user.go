package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Status    int16
	Role      string
	CreatedAt time.Time
}
