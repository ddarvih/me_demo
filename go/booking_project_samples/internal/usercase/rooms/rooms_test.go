package rooms

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/internal/repository"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
)

type stubRooms struct {
	createRoom func(context.Context, string, *string, *int) (entity.Room, error)
}

func (s *stubRooms) CreateRoom(ctx context.Context, name string, desc *string, cap *int) (entity.Room, error) {
	if s.createRoom != nil {
		return s.createRoom(ctx, name, desc, cap)
	}
	return entity.Room{}, nil
}
func (s *stubRooms) ListRooms(context.Context) ([]entity.Room, error)        { panic("not used") }
func (s *stubRooms) GetRoom(context.Context, uuid.UUID) (entity.Room, error) { panic("not used") }
func (s *stubRooms) CreateSchedule(context.Context, uuid.UUID, []int, string, string) (entity.Schedule, error) {
	panic("not used")
}
func (s *stubRooms) GetScheduleByRoom(context.Context, uuid.UUID) (entity.Schedule, bool, error) {
	panic("not used")
}
func (s *stubRooms) UpsertSlots(context.Context, []entity.Slot) error { panic("not used") }
func (s *stubRooms) ListAvailableSlotsForDay(context.Context, uuid.UUID, time.Time, time.Time) ([]entity.Slot, error) {
	panic("not used")
}
func (s *stubRooms) GetSlot(context.Context, uuid.UUID) (entity.Slot, bool, error) { panic("not used") }
func (s *stubRooms) CreateBooking(context.Context, uuid.UUID, uuid.UUID, *string) (entity.Booking, error) {
	panic("not used")
}
func (s *stubRooms) UpdateBookingConferenceLink(context.Context, uuid.UUID, string) (entity.Booking, error) {
	panic("not used")
}
func (s *stubRooms) GetBooking(context.Context, uuid.UUID) (entity.Booking, bool, error) {
	panic("not used")
}
func (s *stubRooms) ListBookingsPaged(context.Context, int, int) ([]entity.Booking, entity.Pagination, error) {
	panic("not used")
}
func (s *stubRooms) ListMyFutureBookings(context.Context, uuid.UUID, time.Time) ([]entity.Booking, error) {
	panic("not used")
}
func (s *stubRooms) CancelBooking(context.Context, uuid.UUID, uuid.UUID) (entity.Booking, error) {
	panic("not used")
}

var _ repository.Store = (*stubRooms)(nil)

func TestCreate_EmptyName(t *testing.T) {
	u := New(&stubRooms{})
	_, err := u.Create(context.Background(), entity.AuthUser{}, "", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if apperrors.AsAppError(err).Code != apperrors.InvalidRequest {
		t.Fatalf("got %v", err)
	}
}

func TestCreate_OK(t *testing.T) {
	rid := uuid.New()
	u := New(&stubRooms{
		createRoom: func(_ context.Context, name string, _ *string, _ *int) (entity.Room, error) {
			return entity.Room{ID: rid, Name: name}, nil
		},
	})
	r, err := u.Create(context.Background(), entity.AuthUser{}, "Hall", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.Name != "Hall" || r.ID != rid {
		t.Fatalf("room %+v", r)
	}
}
