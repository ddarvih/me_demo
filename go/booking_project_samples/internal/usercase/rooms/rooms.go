package rooms

import (
	"context"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/internal/repository"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
)

type UseCase struct {
	store repository.Store
}

func New(store repository.Store) *UseCase {
	return &UseCase{store: store}
}

func (u *UseCase) Create(ctx context.Context, au entity.AuthUser, name string, description *string, capacity *int) (entity.Room, error) {
	if name == "" {
		return entity.Room{}, apperrors.BadRequest(apperrors.InvalidRequest, "name is required")
	}
	return u.store.CreateRoom(ctx, name, description, capacity)
}

func (u *UseCase) List(ctx context.Context, _ entity.AuthUser) ([]entity.Room, error) {
	return u.store.ListRooms(ctx)
}
