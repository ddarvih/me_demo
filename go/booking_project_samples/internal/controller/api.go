package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	authctl "github.com/test-backend-ddarvih/room-booking/internal/controller/auth"
	authuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/auth"

	"github.com/test-backend-ddarvih/room-booking/internal/config"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/handlers/sys"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/middleware"

	bookingctl "github.com/test-backend-ddarvih/room-booking/internal/controller/bookings"
	roomctl "github.com/test-backend-ddarvih/room-booking/internal/controller/rooms"
	schedulectl "github.com/test-backend-ddarvih/room-booking/internal/controller/schedules"
	slotctl "github.com/test-backend-ddarvih/room-booking/internal/controller/slots"

	bookingsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/bookings"
	roomsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/rooms"
	schedulesuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/schedules"
	slotsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/slots"
)

type Deps struct {
	Config    config.Config
	Auth      *authuc.UseCase
	Rooms     *roomsuc.UseCase
	Schedules *schedulesuc.UseCase
	Slots     *slotsuc.UseCase
	Bookings  *bookingsuc.UseCase
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	roomH := roomctl.NewHandler(d.Rooms)
	scheduleH := schedulectl.NewHandler(d.Schedules)
	slotH := slotctl.NewHandler(d.Slots)
	bookingH := bookingctl.NewHandler(d.Bookings)

	r.Get("/_info", sys.Health)
	r.Get("/health", sys.Health)
	r.Post("/dummyLogin", authctl.DummyLogin(d.Auth))
	//r.Post("/dummyLogin", authctl.DummyLogin(d.Config))

	r.Group(func(r chi.Router) {
		r.Use(middleware.BearerJWT(d.Config.JWTSecret))
		r.Group(func(r chi.Router) {
			r.Use(middleware.OnlyAdmin)

			r.Post("/rooms/create", roomH.Create)
			r.Post("/rooms/{roomId}/schedule/create", scheduleH.Create)
			r.Get("/bookings/list", bookingH.ListAll)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.OnlyUser)

			r.Post("/bookings/create", bookingH.Create)
			r.Get("/bookings/my", bookingH.ListMy)
			r.Post("/bookings/{bookingId}/cancel", bookingH.Cancel)
		})

		r.Get("/rooms/list", roomH.List)
		r.Get("/rooms/{roomId}/slots/list", slotH.List)
	})

	r.MethodNotAllowed(sys.NotFoundOrStub)
	r.NotFound(sys.NotFoundOrStub)

	return r
}
