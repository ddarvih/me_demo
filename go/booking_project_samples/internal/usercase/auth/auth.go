package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/test-backend-ddarvih/room-booking/internal/config"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
)

type UseCase struct {
	config config.Config
}

func New(cfg config.Config) *UseCase {
	return &UseCase{config: cfg}
}

func (u *UseCase) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Можно вынести 24ч в конфиг
		"iss":     "room-booking",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.config.JWTSecret))
}

func ValidateDummyRole(role string) error {
	if role != entity.RoleAdmin && role != entity.RoleUser {
		return apperrors.BadRequest(apperrors.InvalidRequest, "role must be admin or user")
	}
	return nil
}

func DummyUserIDForRole(role string) (uuid.UUID, error) {
	switch role {
	case entity.RoleAdmin:
		return entity.AdminUser.ID, nil
	case entity.RoleUser:
		return entity.StandardUser.ID, nil
	default:
		return uuid.Nil, fmt.Errorf("invalid role")
	}
}
