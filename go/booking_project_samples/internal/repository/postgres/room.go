package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
)

const (
	createRoomQuery = `
		INSERT INTO rooms (name, description, capacity)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, capacity, created_at`

	listRoomsQuery = `
		SELECT id, name, description, capacity, created_at 
		FROM rooms 
		ORDER BY created_at`

	getRoomQuery = `
		SELECT id, name, description, capacity, created_at 
		FROM rooms 
		WHERE id = $1`
)

func scanRoom(row pgx.Row) (entity.Room, error) {
	var r entity.Room
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.Capacity, &r.CreatedAt)
	return r, err
}

func (s *Store) CreateRoom(ctx context.Context, name string, description *string, capacity *int) (entity.Room, error) {
	r, err := scanRoom(s.pool.QueryRow(ctx, createRoomQuery, name, description, capacity))
	if err != nil {
		return entity.Room{}, apperrors.Wrap(apperrors.InternalError, 500, "create room", err)
	}
	return r, nil
}

func (s *Store) ListRooms(ctx context.Context) ([]entity.Room, error) {
	rows, err := s.pool.Query(ctx, listRoomsQuery)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.InternalError, 500, "list rooms", err)
	}
	defer rows.Close()

	var list []entity.Room
	for rows.Next() {
		r, err := scanRoom(rows)
		if err != nil {
			return nil, apperrors.Wrap(apperrors.InternalError, 500, "scan room", err)
		}
		list = append(list, r)
	}

	return list, rows.Err()
}

func (s *Store) GetRoom(ctx context.Context, id uuid.UUID) (entity.Room, error) {
	r, err := scanRoom(s.pool.QueryRow(ctx, getRoomQuery, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Room{}, apperrors.NotFoundCode(apperrors.RoomNotFound, "room not found")
	}
	if err != nil {
		return entity.Room{}, apperrors.Wrap(apperrors.InternalError, 500, "get room", err)
	}
	return r, nil
}
