package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
)

type Store interface {
	CreateRoom(ctx context.Context, name string, description *string, capacity *int) (entity.Room, error)
	ListRooms(ctx context.Context) ([]entity.Room, error)
	GetRoom(ctx context.Context, id uuid.UUID) (entity.Room, error)

	CreateSchedule(ctx context.Context, roomID uuid.UUID, days []int, startTime, endTime string) (entity.Schedule, error)
	GetScheduleByRoom(ctx context.Context, roomID uuid.UUID) (entity.Schedule, bool, error)

	UpsertSlots(ctx context.Context, slots []entity.Slot) error
	ListAvailableSlotsForDay(ctx context.Context, roomID uuid.UUID, dayStart, dayEnd time.Time) ([]entity.Slot, error)
	GetSlot(ctx context.Context, id uuid.UUID) (entity.Slot, bool, error)

	CreateBooking(ctx context.Context, slotID, userID uuid.UUID, link *string) (entity.Booking, error)
	UpdateBookingConferenceLink(ctx context.Context, bookingID uuid.UUID, link string) (entity.Booking, error)
	GetBooking(ctx context.Context, id uuid.UUID) (entity.Booking, bool, error)
	ListBookingsPaged(ctx context.Context, page, pageSize int) ([]entity.Booking, entity.Pagination, error)
	ListMyFutureBookings(ctx context.Context, userID uuid.UUID, now time.Time) ([]entity.Booking, error)
	CancelBooking(ctx context.Context, bookingID, userID uuid.UUID) (entity.Booking, error)
}
