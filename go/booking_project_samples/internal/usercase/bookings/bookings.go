package bookings

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/internal/repository"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
)

type ConferenceProvider interface {
	CreateLink(ctx context.Context, bookingID uuid.UUID) (string, error)
}

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type UseCase struct {
	store        repository.Store
	confProvider ConferenceProvider
}

func New(store repository.Store, conf ConferenceProvider) *UseCase {
	return &UseCase{
		store:        store,
		confProvider: conf,
	}
}

func (u *UseCase) Create(ctx context.Context, au entity.AuthUser, slotID uuid.UUID, createLink bool) (entity.Booking, string, error) {
	slot, ok, err := u.store.GetSlot(ctx, slotID)
	if err != nil {
		return entity.Booking{}, "", err
	}
	if !ok {
		return entity.Booking{}, "", apperrors.NotFoundCode(apperrors.SlotNotFound, "slot not found")
	}

	if slot.Start.Before(time.Now().UTC()) {
		return entity.Booking{}, "", apperrors.BadRequest(apperrors.InvalidRequest, "cannot book a slot in the past")
	}

	b, err := u.store.CreateBooking(ctx, slotID, au.ID, nil)
	if err != nil {
		return entity.Booking{}, "", err
	}

	var warning string
	if createLink {
		link, err := u.confProvider.CreateLink(ctx, b.ID)
		if err != nil {
			log.Printf("ERROR: conference service failure: %v", err)
			warning = "booking created, but without requested conference link"
			return b, warning, nil
		}
		updated, err := u.store.UpdateBookingConferenceLink(ctx, b.ID, link)
		if err != nil {
			log.Printf("ERROR: failed to persist conference link after external success: %v", err)
			warning = "booking created; conference link could not be saved"
			return b, warning, nil
		}
		return updated, warning, nil
	}

	return b, warning, nil
}

func (u *UseCase) ListAll(ctx context.Context, _ entity.AuthUser, page, pageSize int) ([]entity.Booking, entity.Pagination, error) {
	if page < 1 {
		page = defaultPage
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		return nil, entity.Pagination{}, apperrors.BadRequest(apperrors.InvalidRequest, "pageSize must be at most 100")
	}
	return u.store.ListBookingsPaged(ctx, page, pageSize)
}

func (u *UseCase) ListMy(ctx context.Context, au entity.AuthUser) ([]entity.Booking, error) {
	return u.store.ListMyFutureBookings(ctx, au.ID, time.Now().UTC())
}

func (u *UseCase) Cancel(ctx context.Context, au entity.AuthUser, bookingID uuid.UUID) (entity.Booking, error) {
	return u.store.CancelBooking(ctx, bookingID, au.ID)
}
