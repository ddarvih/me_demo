package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/test-backend-ddarvih/room-booking/internal/infrastructure/conference"
	authuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/auth"

	"github.com/test-backend-ddarvih/room-booking/internal/config"
	"github.com/test-backend-ddarvih/room-booking/internal/controller"
	"github.com/test-backend-ddarvih/room-booking/internal/repository/postgres"
	bookingsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/bookings"
	roomsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/rooms"
	schedulesuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/schedules"
	slotsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/slots"
)

type App struct {
	Router http.Handler
}

func New(cfg config.Config, pool *pgxpool.Pool) *App {
	store := postgres.NewStore(pool)

	confMock := conference.NewMock()

	authUC := authuc.New(cfg)
	roomsUC := roomsuc.New(store)
	schedulesUC := schedulesuc.New(store)
	slotsUC := slotsuc.New(store)
	bookingsUC := bookingsuc.New(store, confMock)

	r := controller.NewRouter(controller.Deps{
		Config:    cfg,
		Auth:      authUC,
		Rooms:     roomsUC,
		Schedules: schedulesUC,
		Slots:     slotsUC,
		Bookings:  bookingsUC,
	})

	return &App{Router: r}
}
