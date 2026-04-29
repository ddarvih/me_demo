package bookings

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

type stubStore struct {
	getSlot               func(ctx context.Context, id uuid.UUID) (entity.Slot, bool, error)
	createBooking         func(ctx context.Context, slotID, userID uuid.UUID, link *string) (entity.Booking, error)
	updateBookingConfLink func(ctx context.Context, bookingID uuid.UUID, link string) (entity.Booking, error)
}

func (s *stubStore) CreateRoom(context.Context, string, *string, *int) (entity.Room, error) {
	panic("not used")
}
func (s *stubStore) ListRooms(context.Context) ([]entity.Room, error) { panic("not used") }
func (s *stubStore) GetRoom(context.Context, uuid.UUID) (entity.Room, error) {
	panic("not used")
}
func (s *stubStore) CreateSchedule(context.Context, uuid.UUID, []int, string, string) (entity.Schedule, error) {
	panic("not used")
}
func (s *stubStore) GetScheduleByRoom(context.Context, uuid.UUID) (entity.Schedule, bool, error) {
	panic("not used")
}
func (s *stubStore) UpsertSlots(context.Context, []entity.Slot) error { panic("not used") }
func (s *stubStore) ListAvailableSlotsForDay(context.Context, uuid.UUID, time.Time, time.Time) ([]entity.Slot, error) {
	panic("not used")
}
func (s *stubStore) GetSlot(ctx context.Context, id uuid.UUID) (entity.Slot, bool, error) {
	if s.getSlot != nil {
		return s.getSlot(ctx, id)
	}
	return entity.Slot{}, false, nil
}
func (s *stubStore) CreateBooking(ctx context.Context, slotID, userID uuid.UUID, link *string) (entity.Booking, error) {
	if s.createBooking != nil {
		return s.createBooking(ctx, slotID, userID, link)
	}
	return entity.Booking{}, errors.New("createBooking not stubbed")
}
func (s *stubStore) UpdateBookingConferenceLink(ctx context.Context, bookingID uuid.UUID, link string) (entity.Booking, error) {
	if s.updateBookingConfLink != nil {
		return s.updateBookingConfLink(ctx, bookingID, link)
	}
	return entity.Booking{}, errors.New("updateBookingConfLink not stubbed")
}
func (s *stubStore) GetBooking(context.Context, uuid.UUID) (entity.Booking, bool, error) {
	panic("not used")
}
func (s *stubStore) ListBookingsPaged(context.Context, int, int) ([]entity.Booking, entity.Pagination, error) {
	panic("not used")
}
func (s *stubStore) ListMyFutureBookings(context.Context, uuid.UUID, time.Time) ([]entity.Booking, error) {
	panic("not used")
}
func (s *stubStore) CancelBooking(context.Context, uuid.UUID, uuid.UUID) (entity.Booking, error) {
	panic("not used")
}

var _ repository.Store = (*stubStore)(nil)

type stubConf struct {
	link string
	err  error
}

func (c *stubConf) CreateLink(_ context.Context, bookingID uuid.UUID) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if c.link != "" {
		return c.link, nil
	}
	return "https://conf.example/" + bookingID.String()[:8], nil
}

func TestUseCase_Create_PastSlot(t *testing.T) {
	slotID := uuid.New()
	past := time.Now().UTC().Add(-time.Hour)
	u := New(&stubStore{
		getSlot: func(context.Context, uuid.UUID) (entity.Slot, bool, error) {
			return entity.Slot{ID: slotID, Start: past, End: past.Add(30 * time.Minute)}, true, nil
		},
	}, &stubConf{})

	_, _, err := u.Create(context.Background(), entity.AuthUser{ID: uuid.New()}, slotID, false)
	if err == nil {
		t.Fatal("expected error for past slot")
	}
	ae := apperrors.AsAppError(err)
	if ae.Code != apperrors.InvalidRequest {
		t.Fatalf("code: %s", ae.Code)
	}
}

func TestUseCase_Create_SlotNotFound(t *testing.T) {
	u := New(&stubStore{
		getSlot: func(context.Context, uuid.UUID) (entity.Slot, bool, error) {
			return entity.Slot{}, false, nil
		},
	}, &stubConf{})

	_, _, err := u.Create(context.Background(), entity.AuthUser{ID: uuid.New()}, uuid.New(), false)
	if err == nil {
		t.Fatal("expected error")
	}
	if apperrors.AsAppError(err).Code != apperrors.SlotNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestUseCase_Create_WithoutConference(t *testing.T) {
	slotID := uuid.New()
	start := time.Now().UTC().Add(24 * time.Hour)
	bid := uuid.New()
	u := New(&stubStore{
		getSlot: func(context.Context, uuid.UUID) (entity.Slot, bool, error) {
			return entity.Slot{ID: slotID, Start: start, End: start.Add(30 * time.Minute)}, true, nil
		},
		createBooking: func(_ context.Context, sid, uid uuid.UUID, link *string) (entity.Booking, error) {
			if link != nil {
				t.Fatal("expected nil conference link on insert")
			}
			return entity.Booking{ID: bid, SlotID: sid, UserID: uid, Status: entity.BookingActive}, nil
		},
	}, &stubConf{})

	b, warn, err := u.Create(context.Background(), entity.AuthUser{ID: uuid.New()}, slotID, false)
	if err != nil {
		t.Fatal(err)
	}
	if warn != "" {
		t.Fatalf("unexpected warning: %q", warn)
	}
	if b.ID != bid {
		t.Fatalf("booking id")
	}
}

func TestUseCase_Create_ConferenceServiceFails_BookingStillCreated(t *testing.T) {
	slotID := uuid.New()
	start := time.Now().UTC().Add(48 * time.Hour)
	bid := uuid.New()
	u := New(&stubStore{
		getSlot: func(context.Context, uuid.UUID) (entity.Slot, bool, error) {
			return entity.Slot{ID: slotID, Start: start, End: start.Add(30 * time.Minute)}, true, nil
		},
		createBooking: func(context.Context, uuid.UUID, uuid.UUID, *string) (entity.Booking, error) {
			return entity.Booking{ID: bid, Status: entity.BookingActive}, nil
		},
	}, &stubConf{err: errors.New("service unavailable")})

	b, warn, err := u.Create(context.Background(), entity.AuthUser{ID: uuid.New()}, slotID, true)
	if err != nil {
		t.Fatal(err)
	}
	if warn == "" {
		t.Fatal("expected warning when conference fails")
	}
	if b.ConferenceLink != nil {
		t.Fatal("expected no conference link")
	}
}

func TestUseCase_Create_ConferenceSuccess_PersistsLink(t *testing.T) {
	slotID := uuid.New()
	start := time.Now().UTC().Add(72 * time.Hour)
	bid := uuid.New()
	link := "https://conf.example/abc"
	u := New(&stubStore{
		getSlot: func(context.Context, uuid.UUID) (entity.Slot, bool, error) {
			return entity.Slot{ID: slotID, Start: start, End: start.Add(30 * time.Minute)}, true, nil
		},
		createBooking: func(context.Context, uuid.UUID, uuid.UUID, *string) (entity.Booking, error) {
			return entity.Booking{ID: bid, Status: entity.BookingActive}, nil
		},
		updateBookingConfLink: func(_ context.Context, id uuid.UUID, l string) (entity.Booking, error) {
			if id != bid || l != link {
				t.Fatalf("update args: %v %q", id, l)
			}
			return entity.Booking{ID: bid, Status: entity.BookingActive, ConferenceLink: &l}, nil
		},
	}, &stubConf{link: link})

	b, warn, err := u.Create(context.Background(), entity.AuthUser{ID: uuid.New()}, slotID, true)
	if err != nil {
		t.Fatal(err)
	}
	if warn != "" {
		t.Fatalf("unexpected warning: %q", warn)
	}
	if b.ConferenceLink == nil || *b.ConferenceLink != link {
		t.Fatalf("conference link not set: %+v", b.ConferenceLink)
	}
}

func TestUseCase_Create_ConferenceOK_DBUpdateFails(t *testing.T) {
	slotID := uuid.New()
	start := time.Now().UTC().Add(96 * time.Hour)
	bid := uuid.New()
	u := New(&stubStore{
		getSlot: func(context.Context, uuid.UUID) (entity.Slot, bool, error) {
			return entity.Slot{ID: slotID, Start: start, End: start.Add(30 * time.Minute)}, true, nil
		},
		createBooking: func(context.Context, uuid.UUID, uuid.UUID, *string) (entity.Booking, error) {
			return entity.Booking{ID: bid, Status: entity.BookingActive}, nil
		},
		updateBookingConfLink: func(context.Context, uuid.UUID, string) (entity.Booking, error) {
			return entity.Booking{}, errors.New("db error")
		},
	}, &stubConf{link: "https://x"})

	b, warn, err := u.Create(context.Background(), entity.AuthUser{ID: uuid.New()}, slotID, true)
	if err != nil {
		t.Fatal(err)
	}
	if warn == "" {
		t.Fatal("expected warning when DB update fails after conference OK")
	}
	if b.ConferenceLink != nil {
		t.Fatal("booking should be without persisted link")
	}
}

func TestUseCase_ListAll_PageSizeCap(t *testing.T) {
	u := New(&stubStore{}, &stubConf{})
	_, _, err := u.ListAll(context.Background(), entity.AuthUser{}, 1, 200)
	if err == nil {
		t.Fatal("expected error")
	}
	if apperrors.AsAppError(err).Code != apperrors.InvalidRequest {
		t.Fatalf("got %v", err)
	}
}
