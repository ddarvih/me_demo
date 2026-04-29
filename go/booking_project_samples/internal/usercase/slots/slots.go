package slots

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/domain"
	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/internal/repository"
)

type UseCase struct {
	store repository.Store
}

func New(store repository.Store) *UseCase {
	return &UseCase{store: store}
}

func (u *UseCase) ListAvailable(ctx context.Context, _ entity.AuthUser, roomID uuid.UUID, dateUTC time.Time) ([]entity.Slot, error) {
	if _, err := u.store.GetRoom(ctx, roomID); err != nil {
		return nil, err
	}
	sch, ok, err := u.store.GetScheduleByRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return []entity.Slot{}, nil
	}
	generated := domain.GenerateDaySlots(roomID, dateUTC, sch)
	if len(generated) == 0 {
		return []entity.Slot{}, nil
	}
	if err := u.store.UpsertSlots(ctx, generated); err != nil {
		return nil, err
	}
	dayStart := time.Date(dateUTC.Year(), dateUTC.Month(), dateUTC.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)
	return u.store.ListAvailableSlotsForDay(ctx, roomID, dayStart, dayEnd)
}
