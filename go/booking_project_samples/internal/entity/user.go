package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type AuthUser struct {
	ID   uuid.UUID
	Role string
}

type User struct {
	ID        uuid.UUID
	Email     string
	Role      string
	CreatedAt time.Time
}

var (
	AdminUser = User{
		ID:    uuid.MustParse("a0000001-0000-4000-8000-000000000001"),
		Role:  RoleAdmin,
		Email: "admin@example.com",
	}
	StandardUser = User{
		ID:    uuid.MustParse("a0000001-0000-4000-8000-000000000002"),
		Role:  RoleUser,
		Email: "user@example.com",
	}
)
