package auth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/test-backend-ddarvih/room-booking/internal/config"
	"github.com/test-backend-ddarvih/room-booking/internal/entity"
)

func TestAuth_GenerateToken(t *testing.T) {
	cfg := config.Config{JWTSecret: "test-secret"}
	uc := New(cfg)

	t.Run("Generate for Admin", func(t *testing.T) {
		id := uuid.New()
		token, err := uc.GenerateAccessToken(id, entity.RoleAdmin)
		if err != nil {
			t.Errorf("failed to generate token: %v", err)
		}
		if token == "" {
			t.Error("token is empty")
		}
	})
}

func TestAuth_DummyRoles(t *testing.T) {
	t.Run("Roles correct validation", func(t *testing.T) {
		if err := ValidateDummyRole("admin"); err != nil {
			t.Errorf("expected admin to be valid: %v", err)
		}
		if err := ValidateDummyRole("user"); err != nil {
			t.Errorf("expected user to be valid: %v", err)
		}
	})

	t.Run("Invalid role error", func(t *testing.T) {
		if err := ValidateDummyRole("hacker"); err == nil {
			t.Error("expected invalid role to return error")
		}
	})

	t.Run("Static IDs check", func(t *testing.T) {
		adminID, _ := DummyUserIDForRole("admin")
		userID, _ := DummyUserIDForRole("user")
		if adminID == userID {
			t.Error("expected admin and user to have different static IDs")
		}
	})
}
