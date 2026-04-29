package bookings

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/converter"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/middleware"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	au, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("unauthorized"))
		return
	}
	bid, err := uuid.Parse(chi.URLParam(r, "bookingId"))
	if err != nil {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "invalid booking id"))
		return
	}
	b, err := h.UC.Cancel(r.Context(), au, bid)
	if err != nil {
		jsonresp.WriteAppError(w, err)
		return
	}
	jsonresp.WriteJSON(w, http.StatusOK, openapi.CancelBookingResponse{Booking: converter.BookingToOAPI(b)})
}
