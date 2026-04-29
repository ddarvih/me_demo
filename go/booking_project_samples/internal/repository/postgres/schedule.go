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
	createScheduleQuery = `
		INSERT INTO schedules (room_id, days_of_week, start_time, end_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id, room_id, days_of_week, start_time, end_time, created_at`

	getScheduleByRoomQuery = `
		SELECT id, room_id, days_of_week, start_time, end_time, created_at 
		FROM schedules 
		WHERE room_id = $1`
)

func scanSchedule(row pgx.Row) (entity.Schedule, error) {
	var sch entity.Schedule
	var daysDB []int32

	err := row.Scan(&sch.ID, &sch.RoomID, &daysDB, &sch.StartTime, &sch.EndTime, &sch.CreatedAt)
	if err != nil {
		return entity.Schedule{}, err
	}

	sch.DaysOfWeek = int32sToInts(daysDB)
	return sch, nil
}

func (s *Store) CreateSchedule(ctx context.Context, roomID uuid.UUID, days []int, startTime, endTime string) (entity.Schedule, error) {
	sch, err := scanSchedule(s.pool.QueryRow(ctx, createScheduleQuery, roomID, days, startTime, endTime))

	if isUniqueViolation(err) {
		return entity.Schedule{}, apperrors.Conflict(apperrors.ScheduleExists, "schedule already exists")
	}
	if err != nil {
		return entity.Schedule{}, apperrors.Wrap(apperrors.InternalError, 500, "create schedule", err)
	}

	return sch, nil
}

func (s *Store) GetScheduleByRoom(ctx context.Context, roomID uuid.UUID) (entity.Schedule, bool, error) {
	sch, err := scanSchedule(s.pool.QueryRow(ctx, getScheduleByRoomQuery, roomID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Schedule{}, false, nil
	}
	if err != nil {
		return entity.Schedule{}, false, apperrors.Wrap(apperrors.InternalError, 500, "get schedule", err)
	}
	return sch, true, nil
}
