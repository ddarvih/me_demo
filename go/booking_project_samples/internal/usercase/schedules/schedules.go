package schedules

import (
	"context"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/domain"
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

func (u *UseCase) Create(ctx context.Context, _ entity.AuthUser, roomID uuid.UUID, days []int, startTime, endTime string) (entity.Schedule, error) {
	if _, err := u.store.GetRoom(ctx, roomID); err != nil {
		return entity.Schedule{}, err
	}
	if _, ok, _ := u.store.GetScheduleByRoom(ctx, roomID); ok {
		return entity.Schedule{}, apperrors.Conflict(apperrors.ScheduleExists, "schedule for this room already exists and cannot be changed")
	}
	if err := validateScheduleDays(days); err != nil {
		return entity.Schedule{}, err
	}

	sh, sm, okStart := domain.ParseHHMM(startTime)
	eh, em, okEnd := domain.ParseHHMM(endTime)
	if !okStart || !okEnd {
		return entity.Schedule{}, apperrors.BadRequest(apperrors.InvalidRequest, "invalid time format, expected HH:MM")
	}
	if sh > eh || (sh == eh && sm >= em) {
		return entity.Schedule{}, apperrors.BadRequest(apperrors.InvalidRequest, "startTime must be before endTime")
	}
	return u.store.CreateSchedule(ctx, roomID, days, startTime, endTime)
}

func validateScheduleDays(days []int) error {
	if len(days) == 0 {
		return apperrors.BadRequest(apperrors.InvalidRequest, "daysOfWeek is required")
	}
	seen := make(map[int]struct{})
	for _, d := range days {
		if d < 1 || d > 7 {
			return apperrors.BadRequest(apperrors.InvalidRequest, "daysOfWeek values must be between 1 and 7")
		}
		if _, ok := seen[d]; ok {
			continue
		}
		seen[d] = struct{}{}
	}
	return nil
}
