package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/entity"
)

var slotNamespace = uuid.MustParse("018e1c50-0000-7000-8000-000000000001")

func APIDayOfWeek(t time.Time) int {
	w := t.UTC().Weekday()
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

func containsAPIDay(days []int, t time.Time) bool {
	d := APIDayOfWeek(t)
	for _, x := range days {
		if x == d {
			return true
		}
	}
	return false
}

func SlotIDFor(roomID uuid.UUID, start time.Time) uuid.UUID {
	payload := roomID.String() + "|" + start.UTC().Format(time.RFC3339Nano)
	return uuid.NewSHA1(slotNamespace, []byte(payload))
}

func ParseHHMM(s string) (h, m int, ok bool) {
	if len(s) != 5 || s[2] != ':' {
		return 0, 0, false
	}
	for i, r := range s {
		if i == 2 {
			continue
		}
		if r < '0' || r > '9' {
			return 0, 0, false
		}
	}
	h = int(s[0]-'0')*10 + int(s[1]-'0')
	m = int(s[3]-'0')*10 + int(s[4]-'0')
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, false
	}
	return h, m, true
}

func GenerateDaySlots(roomID uuid.UUID, dateUTC time.Time, sch entity.Schedule) []entity.Slot {
	dateUTC = time.Date(dateUTC.Year(), dateUTC.Month(), dateUTC.Day(), 0, 0, 0, 0, time.UTC)
	if !containsAPIDay(sch.DaysOfWeek, dateUTC) {
		return nil
	}

	sh, sm, ok := ParseHHMM(sch.StartTime)
	if !ok {
		return nil
	}
	eh, em, ok := ParseHHMM(sch.EndTime)
	if !ok {
		return nil
	}
	if sh > eh || (sh == eh && sm >= em) {
		return nil
	}

	dayStart := dateUTC.Add(time.Duration(sh)*time.Hour + time.Duration(sm)*time.Minute)
	dayEnd := dateUTC.Add(time.Duration(eh)*time.Hour + time.Duration(em)*time.Minute)

	var slots []entity.Slot
	for cur := dayStart; cur.Add(30*time.Minute).Compare(dayEnd) <= 0; cur = cur.Add(30 * time.Minute) {
		end := cur.Add(30 * time.Minute)
		slots = append(slots, entity.Slot{
			ID:     SlotIDFor(roomID, cur),
			RoomID: roomID,
			Start:  cur,
			End:    end,
		})
	}
	return slots
}
