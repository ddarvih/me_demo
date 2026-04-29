package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/test-backend-ddarvih/room-booking/internal/config"
	"github.com/test-backend-ddarvih/room-booking/internal/controller"
	"github.com/test-backend-ddarvih/room-booking/internal/infrastructure/conference"
	authuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/auth"
	bookingsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/bookings"
	roomsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/rooms"
	schedulesuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/schedules"
	slotsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/slots"
)

func TestRouter_InfoReturns200(t *testing.T) {
	cfg := config.Config{JWTSecret: "x"}
	d := controller.Deps{
		Config:    cfg,
		Auth:      authuc.New(cfg),
		Rooms:     roomsuc.New(nil),
		Schedules: schedulesuc.New(nil),
		Slots:     slotsuc.New(nil),
		Bookings:  bookingsuc.New(nil, conference.NewMock()),
	}
	h := controller.NewRouter(d)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/_info", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
}
