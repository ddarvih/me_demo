package schedules

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/internal/repository"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
)

type stubSched struct {
	getRoom           func(context.Context, uuid.UUID) (entity.Room, error)
	getScheduleByRoom func(context.Context, uuid.UUID) (entity.Schedule, bool, error)
	createSchedule    func(context.Context, uuid.UUID, []int, string, string) (entity.Schedule, error)
}

func (s *stubSched) CreateRoom(context.Context, string, *string, *int) (entity.Room, error) {
	panic("not used")
}
func (s *stubSched) ListRooms(context.Context) ([]entity.Room, error) { panic("not used") }
func (s *stubSched) GetRoom(ctx context.Context, id uuid.UUID) (entity.Room, error) {
	if s.getRoom != nil {
		return s.getRoom(ctx, id)
	}
	return entity.Room{}, errors.New("not found")
}
func (s *stubSched) CreateSchedule(ctx context.Context, roomID uuid.UUID, days []int, st, et string) (entity.Schedule, error) {
	if s.createSchedule != nil {
		return s.createSchedule(ctx, roomID, days, st, et)
	}
	panic("not stubbed")
}
func (s *stubSched) GetScheduleByRoom(ctx context.Context, roomID uuid.UUID) (entity.Schedule, bool, error) {
	if s.getScheduleByRoom != nil {
		return s.getScheduleByRoom(ctx, roomID)
	}
	return entity.Schedule{}, false, nil
}
func (s *stubSched) UpsertSlots(context.Context, []entity.Slot) error { panic("not used") }
func (s *stubSched) ListAvailableSlotsForDay(context.Context, uuid.UUID, time.Time, time.Time) ([]entity.Slot, error) {
	panic("not used")
}
func (s *stubSched) GetSlot(context.Context, uuid.UUID) (entity.Slot, bool, error) { panic("not used") }
func (s *stubSched) CreateBooking(context.Context, uuid.UUID, uuid.UUID, *string) (entity.Booking, error) {
	panic("not used")
}
func (s *stubSched) UpdateBookingConferenceLink(context.Context, uuid.UUID, string) (entity.Booking, error) {
	panic("not used")
}
func (s *stubSched) GetBooking(context.Context, uuid.UUID) (entity.Booking, bool, error) {
	panic("not used")
}
func (s *stubSched) ListBookingsPaged(context.Context, int, int) ([]entity.Booking, entity.Pagination, error) {
	panic("not used")
}
func (s *stubSched) ListMyFutureBookings(context.Context, uuid.UUID, time.Time) ([]entity.Booking, error) {
	panic("not used")
}
func (s *stubSched) CancelBooking(context.Context, uuid.UUID, uuid.UUID) (entity.Booking, error) {
	panic("not used")
}

var _ repository.Store = (*stubSched)(nil)

func TestCreate_DuplicateSchedule(t *testing.T) {
	rid := uuid.New()
	u := New(&stubSched{
		getRoom: func(context.Context, uuid.UUID) (entity.Room, error) {
			return entity.Room{ID: rid}, nil
		},
		getScheduleByRoom: func(context.Context, uuid.UUID) (entity.Schedule, bool, error) {
			return entity.Schedule{}, true, nil
		},
	})

	_, err := u.Create(context.Background(), entity.AuthUser{}, rid, []int{1}, "09:00", "18:00")
	if err == nil {
		t.Fatal("expected conflict")
	}
	if apperrors.AsAppError(err).Code != apperrors.ScheduleExists {
		t.Fatalf("got %v", err)
	}
}

func TestCreate_InvalidDay(t *testing.T) {
	rid := uuid.New()
	u := New(&stubSched{
		getRoom: func(context.Context, uuid.UUID) (entity.Room, error) {
			return entity.Room{ID: rid}, nil
		},
		getScheduleByRoom: func(context.Context, uuid.UUID) (entity.Schedule, bool, error) {
			return entity.Schedule{}, false, nil
		},
	})

	_, err := u.Create(context.Background(), entity.AuthUser{}, rid, []int{8}, "09:00", "18:00")
	if err == nil {
		t.Fatal("expected error")
	}
	if apperrors.AsAppError(err).Code != apperrors.InvalidRequest {
		t.Fatalf("got %v", err)
	}
}

func TestCreate_InvalidTimeRange(t *testing.T) {
	rid := uuid.New()
	u := New(&stubSched{
		getRoom: func(context.Context, uuid.UUID) (entity.Room, error) {
			return entity.Room{ID: rid}, nil
		},
		getScheduleByRoom: func(context.Context, uuid.UUID) (entity.Schedule, bool, error) {
			return entity.Schedule{}, false, nil
		},
	})

	_, err := u.Create(context.Background(), entity.AuthUser{}, rid, []int{1}, "18:00", "09:00")
	if err == nil {
		t.Fatal("expected error")
	}
}
