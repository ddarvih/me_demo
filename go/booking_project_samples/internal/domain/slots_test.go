package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
)

func TestAPIDayOfWeek(t *testing.T) {
	day := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	if got := APIDayOfWeek(day); got != 1 {
		t.Fatalf("Mon: got %d want 1", got)
	}
	sun := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)
	if got := APIDayOfWeek(sun); got != 7 {
		t.Fatalf("Sun: got %d want 7", got)
	}
}

func TestGenerateDaySlots(t *testing.T) {
	rid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	sch := entity.Schedule{
		DaysOfWeek: []int{1},
		StartTime:  "09:00",
		EndTime:    "10:30",
	}
	date := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	slots := GenerateDaySlots(rid, date, sch)
	if len(slots) != 3 {
		t.Fatalf("want 3 slots, got %d", len(slots))
	}
	if !slots[0].Start.Equal(time.Date(2026, 6, 10, 9, 0, 0, 0, time.UTC)) {
		t.Fatalf("first start: %v", slots[0].Start)
	}
	if !slots[2].End.Equal(time.Date(2026, 6, 10, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("last end: %v", slots[2].End)
	}
	for _, s := range slots {
		if s.ID != SlotIDFor(rid, s.Start) {
			t.Fatalf("id mismatch for %v", s.Start)
		}
	}
}

func TestGenerateDaySlotsWrongWeekday(t *testing.T) {
	rid := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	sch := entity.Schedule{DaysOfWeek: []int{1}, StartTime: "09:00", EndTime: "10:00"}
	date := time.Date(2024, 6, 11, 0, 0, 0, 0, time.UTC)
	if slots := GenerateDaySlots(rid, date, sch); len(slots) != 0 {
		t.Fatalf("expected no slots but got %d", len(slots))
	}
}
