package slots

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/internal/repository"
)

type stubSlots struct {
	getRoom           func(context.Context, uuid.UUID) (entity.Room, error)
	getScheduleByRoom func(context.Context, uuid.UUID) (entity.Schedule, bool, error)
	upsertSlots       func(context.Context, []entity.Slot) error
	listDay           func(context.Context, uuid.UUID, time.Time, time.Time) ([]entity.Slot, error)
}

func (s *stubSlots) CreateRoom(context.Context, string, *string, *int) (entity.Room, error) {
	panic("nu")
}
func (s *stubSlots) ListRooms(context.Context) ([]entity.Room, error) { panic("nu") }
func (s *stubSlots) GetRoom(ctx context.Context, id uuid.UUID) (entity.Room, error) {
	if s.getRoom != nil {
		return s.getRoom(ctx, id)
	}
	return entity.Room{}, errors.New("nf")
}
func (s *stubSlots) CreateSchedule(context.Context, uuid.UUID, []int, string, string) (entity.Schedule, error) {
	panic("nu")
}
func (s *stubSlots) GetScheduleByRoom(ctx context.Context, roomID uuid.UUID) (entity.Schedule, bool, error) {
	if s.getScheduleByRoom != nil {
		return s.getScheduleByRoom(ctx, roomID)
	}
	return entity.Schedule{}, false, nil
}
func (s *stubSlots) UpsertSlots(ctx context.Context, slots []entity.Slot) error {
	if s.upsertSlots != nil {
		return s.upsertSlots(ctx, slots)
	}
	return nil
}
func (s *stubSlots) ListAvailableSlotsForDay(ctx context.Context, roomID uuid.UUID, a, b time.Time) ([]entity.Slot, error) {
	if s.listDay != nil {
		return s.listDay(ctx, roomID, a, b)
	}
	return nil, nil
}
func (s *stubSlots) GetSlot(context.Context, uuid.UUID) (entity.Slot, bool, error) { panic("nu") }
func (s *stubSlots) CreateBooking(context.Context, uuid.UUID, uuid.UUID, *string) (entity.Booking, error) {
	panic("nu")
}
func (s *stubSlots) UpdateBookingConferenceLink(context.Context, uuid.UUID, string) (entity.Booking, error) {
	panic("nu")
}
func (s *stubSlots) GetBooking(context.Context, uuid.UUID) (entity.Booking, bool, error) { panic("nu") }
func (s *stubSlots) ListBookingsPaged(context.Context, int, int) ([]entity.Booking, entity.Pagination, error) {
	panic("nu")
}
func (s *stubSlots) ListMyFutureBookings(context.Context, uuid.UUID, time.Time) ([]entity.Booking, error) {
	panic("nu")
}
func (s *stubSlots) CancelBooking(context.Context, uuid.UUID, uuid.UUID) (entity.Booking, error) {
	panic("nu")
}

var _ repository.Store = (*stubSlots)(nil)

func TestListAvailable_NoSchedule_Empty(t *testing.T) {
	rid := uuid.New()
	u := New(&stubSlots{
		getRoom: func(context.Context, uuid.UUID) (entity.Room, error) {
			return entity.Room{ID: rid}, nil
		},
		getScheduleByRoom: func(context.Context, uuid.UUID) (entity.Schedule, bool, error) {
			return entity.Schedule{}, false, nil
		},
	})
	day := time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)
	got, err := u.ListAvailable(context.Background(), entity.AuthUser{}, rid, day)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("want empty, got %d", len(got))
	}
}
