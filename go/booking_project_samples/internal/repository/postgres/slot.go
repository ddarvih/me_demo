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
	upsertSlotQuery = `
		INSERT INTO slots (id, room_id, start_ts, end_ts) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (room_id, start_ts) DO NOTHING`

	listAvailableSlotsQuery = `
		SELECT s.id, s.room_id, s.start_ts, s.end_ts
		FROM slots s
		LEFT JOIN bookings b ON b.slot_id = s.id AND b.status = $4
		WHERE s.room_id = $1 AND s.start_ts >= $2 AND s.start_ts < $3 AND b.id IS NULL
		ORDER BY s.start_ts`

	getSlotQuery = `
		SELECT id, room_id, start_ts, end_ts 
		FROM slots 
		WHERE id = $1`
)

func scanSlot(row pgx.Row) (entity.Slot, error) {
	var sl entity.Slot
	err := row.Scan(&sl.ID, &sl.RoomID, &sl.Start, &sl.End)
	return sl, err
}

func (s *Store) UpsertSlots(ctx context.Context, slots []entity.Slot) error {
	batch := &pgx.Batch{}
	for _, sl := range slots {
		batch.Queue(upsertSlotQuery, sl.ID, sl.RoomID, sl.Start, sl.End)
	}
	br := s.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range slots {
		if _, err := br.Exec(); err != nil {
			return apperrors.Wrap(apperrors.InternalError, 500, "upsert slot", err)
		}
	}
	return nil
}

func (s *Store) ListAvailableSlotsForDay(ctx context.Context, roomID uuid.UUID, dayStart, dayEnd time.Time) ([]entity.Slot, error) {
	rows, err := s.pool.Query(ctx, listAvailableSlotsQuery, roomID, dayStart, dayEnd, entity.BookingActive)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.InternalError, 500, "list slots", err)
	}
	defer rows.Close()

	var list []entity.Slot
	for rows.Next() {
		sl, err := scanSlot(rows)
		if err != nil {
			return nil, apperrors.Wrap(apperrors.InternalError, 500, "scan slot", err)
		}
		list = append(list, sl)
	}
	return list, rows.Err()
}

func (s *Store) GetSlot(ctx context.Context, id uuid.UUID) (entity.Slot, bool, error) {
	var sl entity.Slot
	err := s.pool.QueryRow(ctx, getSlotQuery, id).Scan(&sl.ID, &sl.RoomID, &sl.Start, &sl.End)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Slot{}, false, nil
	}
	if err != nil {
		return entity.Slot{}, false, apperrors.Wrap(apperrors.InternalError, 500, "get slot", err)
	}
	return sl, true, nil
}
