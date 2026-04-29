package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

const baseURL = "http://localhost:8080"

func skipIfNoE2E(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_E2E") == "" {
		t.Skip("set RUN_E2E=1 on localhost:8080")
	}
}

func TestIntegration_RoomScheduleBooking(t *testing.T) {
	skipIfNoE2E(t)
	client := &http.Client{Timeout: 15 * time.Second}
	adminToken := login(t, client, "admin")
	roomName := "Meeting Room " + uuid.New().String()[:8]
	roomID := createRoom(t, client, adminToken, roomName)
	createSchedule(t, client, adminToken, roomID)
	userToken := login(t, client, "user")
	tomorrow := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02")
	slotID := getFirstSlot(t, client, userToken, roomID, tomorrow)
	bookingID := createBooking(t, client, userToken, slotID, false)
	if bookingID == uuid.Nil {
		t.Fatal("nil booking id")
	}
}

func TestIntegration_CancelBooking(t *testing.T) {
	skipIfNoE2E(t)
	client := &http.Client{Timeout: 15 * time.Second}
	adminToken := login(t, client, "admin")
	roomID := createRoom(t, client, adminToken, "Cancel Test "+uuid.New().String()[:8])
	createSchedule(t, client, adminToken, roomID)
	userToken := login(t, client, "user")
	tomorrow := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02")
	slotID := getFirstSlot(t, client, userToken, roomID, tomorrow)
	bookingID := createBooking(t, client, userToken, slotID, false)
	cancelBooking(t, client, userToken, bookingID)
	cancelBooking(t, client, userToken, bookingID)
}

func login(t *testing.T, c *http.Client, role string) string {
	t.Helper()
	body, _ := json.Marshal(openapi.DummyLoginJSONRequestBody{Role: role})
	resp, err := c.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Login request failed for %s: %v", role, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed for %s: %d", role, resp.StatusCode)
	}
	var r openapi.Token
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("decode token: %v", err)
	}
	if r.Token == "" {
		t.Fatalf("empty token for %s", role)
	}
	return r.Token
}

func createRoom(t *testing.T, c *http.Client, token, name string) uuid.UUID {
	t.Helper()
	capVal := 10
	body, _ := json.Marshal(openapi.CreateRoomJSONRequestBody{
		Name:     name,
		Capacity: &capVal,
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/rooms/create", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Create room request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create room fail: %d", resp.StatusCode)
	}
	var r openapi.CreateRoomResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("decode room: %v", err)
	}
	return r.Room.Id
}

func createSchedule(t *testing.T, c *http.Client, token string, roomID uuid.UUID) {
	t.Helper()
	url := fmt.Sprintf("%s/rooms/%s/schedule/create", baseURL, roomID.String())
	body, _ := json.Marshal(openapi.Schedule{
		RoomId:     roomID,
		DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7},
		StartTime:  "09:00",
		EndTime:    "18:00",
	})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Create Schedule request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create Schedule fail: %d", resp.StatusCode)
	}
}

func getFirstSlot(t *testing.T, c *http.Client, token string, roomID uuid.UUID, date string) uuid.UUID {
	t.Helper()
	url := fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID.String(), date)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Get Slots request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get Slots fail: %d", resp.StatusCode)
	}
	var r openapi.ListSlotsResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("decode slots: %v", err)
	}
	if len(r.Slots) == 0 {
		t.Fatal("No slots")
	}
	return r.Slots[0].Id
}

func createBooking(t *testing.T, c *http.Client, token string, slotID uuid.UUID, withConf bool) uuid.UUID {
	t.Helper()
	body, _ := json.Marshal(openapi.CreateBookingJSONRequestBody{
		SlotId:               slotID,
		CreateConferenceLink: withConf,
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/bookings/create", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Create Booking request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create Booking fail: %d", resp.StatusCode)
	}
	var r openapi.CreateBookingResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("decode booking: %v", err)
	}
	return r.Booking.Id
}

func cancelBooking(t *testing.T, c *http.Client, token string, bookingID uuid.UUID) {
	t.Helper()
	url := fmt.Sprintf("%s/bookings/%s/cancel", baseURL, bookingID.String())
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Cancel Booking request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Cancel Booking fail: %d", resp.StatusCode)
	}
}
