package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
)

const (
	createBookingQuery = `
		INSERT INTO bookings (slot_id, user_id, status, conference_link)
		VALUES ($1, $2, $3, $4)
		RETURNING id, slot_id, user_id, status, conference_link, created_at`

	getBookingQuery = `
		SELECT id, slot_id, user_id, status, conference_link, created_at 
		FROM bookings WHERE id = $1`

	lockBookingQuery = `
		SELECT id, slot_id, user_id, status, conference_link, created_at 
		FROM bookings WHERE id = $1 FOR UPDATE`

	countBookingsQuery = `SELECT COUNT(*) FROM bookings`

	listBookingsPagedQuery = `
		SELECT id, slot_id, user_id, status, conference_link, created_at
		FROM bookings 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	listMyFutureBookingsQuery = `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at
		FROM bookings b
		INNER JOIN slots s ON s.id = b.slot_id
		WHERE b.user_id = $1 AND b.status = $2 AND s.start_ts >= $3
		ORDER BY s.start_ts`

	updateBookingStatusQuery = `
		UPDATE bookings SET status = $1 WHERE id = $2 RETURNING status`

	updateConferenceLinkQuery = `
		UPDATE bookings SET conference_link = $2 WHERE id = $1
		RETURNING id, slot_id, user_id, status, conference_link, created_at`

	checkActiveBookingQuery = `
		SELECT EXISTS(SELECT 1 FROM bookings WHERE slot_id = $1 AND status = $2)`

	lockSlotQuery = `
		SELECT id FROM slots WHERE id = $1 FOR UPDATE`
)

func scanBooking(row pgx.Row) (entity.Booking, error) {
	var b entity.Booking
	var status string
	err := row.Scan(&b.ID, &b.SlotID, &b.UserID, &status, &b.ConferenceLink, &b.CreatedAt)
	if err != nil {
		return entity.Booking{}, err
	}
	b.Status = entity.BookingStatus(status)
	return b, err
}

func (s *Store) CreateBooking(ctx context.Context, slotID, userID uuid.UUID, conferenceLink *string) (entity.Booking, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return entity.Booking{}, apperrors.Internal("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	var dummy uuid.UUID
	err = tx.QueryRow(ctx, lockSlotQuery, slotID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Booking{}, apperrors.NotFoundCode(apperrors.SlotNotFound, "slot not found")
		}
		return entity.Booking{}, apperrors.Internal("failed to lock slot")
	}

	var alreadyOccupied bool
	err = tx.QueryRow(ctx, checkActiveBookingQuery, slotID, entity.BookingActive).Scan(&alreadyOccupied)
	if err != nil {
		return entity.Booking{}, apperrors.Internal("failed to check slot availability")
	}
	if alreadyOccupied {
		return entity.Booking{}, apperrors.Conflict(apperrors.SlotAlreadyBooked, "slot is already occupied by an active booking")
	}

	b, err := scanBooking(tx.QueryRow(ctx, createBookingQuery, slotID, userID, entity.BookingActive, conferenceLink))
	if err != nil {
		if isUniqueViolation(err) {
			return entity.Booking{}, apperrors.Conflict(apperrors.SlotAlreadyBooked, "race condition detected: slot was just booked")
		}
		return entity.Booking{}, apperrors.Internal("failed to insert booking")
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.Booking{}, apperrors.Internal("failed to commit booking transaction")
	}

	return b, nil
}

func (s *Store) UpdateBookingConferenceLink(ctx context.Context, bookingID uuid.UUID, link string) (entity.Booking, error) {
	b, err := scanBooking(s.pool.QueryRow(ctx, updateConferenceLinkQuery, bookingID, link))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Booking{}, apperrors.NotFoundCode(apperrors.BookingNotFound, "booking not found")
	}
	if err != nil {
		return entity.Booking{}, apperrors.Wrap(apperrors.InternalError, 500, "update conference link", err)
	}
	return b, nil
}

func (s *Store) GetBooking(ctx context.Context, id uuid.UUID) (entity.Booking, bool, error) {
	b, err := scanBooking(s.pool.QueryRow(ctx, getBookingQuery, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Booking{}, false, nil
	}
	if err != nil {
		return entity.Booking{}, false, apperrors.Wrap(apperrors.InternalError, 500, "get booking", err)
	}
	return b, true, nil
}

func (s *Store) ListBookingsPaged(ctx context.Context, page, pageSize int) ([]entity.Booking, entity.Pagination, error) {
	var total int
	if err := s.pool.QueryRow(ctx, countBookingsQuery).Scan(&total); err != nil {
		return nil, entity.Pagination{}, apperrors.Wrap(apperrors.InternalError, 500, "count bookings", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.pool.Query(ctx, listBookingsPagedQuery, pageSize, offset)
	if err != nil {
		return nil, entity.Pagination{}, apperrors.Wrap(apperrors.InternalError, 500, "list bookings", err)
	}
	defer rows.Close()

	var list []entity.Booking
	for rows.Next() {
		b, err := scanBooking(rows)
		if err != nil {
			return nil, entity.Pagination{}, apperrors.Wrap(apperrors.InternalError, 500, "scan booking", err)
		}
		list = append(list, b)
	}

	return list, entity.Pagination{Page: page, PageSize: pageSize, Total: total}, rows.Err()
}

func (s *Store) ListMyFutureBookings(ctx context.Context, userID uuid.UUID, now time.Time) ([]entity.Booking, error) {
	rows, err := s.pool.Query(ctx, listMyFutureBookingsQuery, userID, entity.BookingActive, now)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.InternalError, 500, "list my bookings", err)
	}
	defer rows.Close()

	var list []entity.Booking
	for rows.Next() {
		b, err := scanBooking(rows)
		if err != nil {
			return nil, apperrors.Wrap(apperrors.InternalError, 500, "scan booking", err)
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (s *Store) CancelBooking(ctx context.Context, bookingID, userID uuid.UUID) (entity.Booking, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return entity.Booking{}, apperrors.Wrap(apperrors.InternalError, 500, "begin tx", err)
	}
	defer tx.Rollback(ctx)

	b, err := scanBooking(tx.QueryRow(ctx, lockBookingQuery, bookingID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Booking{}, apperrors.NotFoundCode(apperrors.BookingNotFound, "booking not found")
	}
	if err != nil {
		return entity.Booking{}, apperrors.Wrap(apperrors.InternalError, 500, "get booking", err)
	}

	if b.UserID != userID {
		return entity.Booking{}, apperrors.ForbiddenMsg("cannot cancel another user's booking")
	}

	if b.Status == entity.BookingCancelled {
		return b, nil
	}

	var newStatus string
	err = tx.QueryRow(ctx, updateBookingStatusQuery, entity.BookingCancelled, bookingID).Scan(&newStatus)
	if err != nil {
		return entity.Booking{}, apperrors.Wrap(apperrors.InternalError, 500, "cancel booking", err)
	}

	b.Status = entity.BookingStatus(newStatus)
	if err := tx.Commit(ctx); err != nil {
		return entity.Booking{}, apperrors.Wrap(apperrors.InternalError, 500, "commit", err)
	}
	return b, nil
}
