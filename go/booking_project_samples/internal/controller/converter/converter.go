package converter

import (
	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

func RoomToOAPI(e entity.Room) openapi.Room {
	r := openapi.Room{Id: e.ID, Name: e.Name}
	r.Description = e.Description
	r.Capacity = e.Capacity
	if !e.CreatedAt.IsZero() {
		t := e.CreatedAt.UTC()
		r.CreatedAt = &t
	}
	return r
}

func ScheduleToOAPI(e entity.Schedule) openapi.Schedule {
	s := openapi.Schedule{
		RoomId:     e.RoomID,
		DaysOfWeek: e.DaysOfWeek,
		StartTime:  e.StartTime,
		EndTime:    e.EndTime,
	}
	s.Id = &e.ID
	if !e.CreatedAt.IsZero() {
		t := e.CreatedAt.UTC()
		s.CreatedAt = &t
	}
	return s
}

func SlotToOAPI(e entity.Slot) openapi.Slot {
	return openapi.Slot{
		Id:     e.ID,
		RoomId: e.RoomID,
		Start:  e.Start.UTC(),
		End:    e.End.UTC(),
	}
}

func BookingToOAPI(e entity.Booking) openapi.Booking {
	b := openapi.Booking{
		Id:             e.ID,
		SlotId:         e.SlotID,
		UserId:         e.UserID,
		Status:         string(e.Status),
		ConferenceLink: e.ConferenceLink,
	}
	if !e.CreatedAt.IsZero() {
		t := e.CreatedAt.UTC()
		b.CreatedAt = &t
	}
	return b
}

func PaginationToOpenAPI(page, pageSize, total int) openapi.Pagination {
	return openapi.Pagination{Page: page, PageSize: pageSize, Total: total}
}
