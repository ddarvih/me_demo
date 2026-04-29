package bookings

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/converter"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/middleware"
	bookingsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/bookings"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

type Handler struct {
	UC *bookingsuc.UseCase
}

func NewHandler(uc *bookingsuc.UseCase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	au, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("unauthorized"))
		return
	}
	var body openapi.CreateBookingJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "invalid json"))
		return
	}
	if body.SlotId == uuid.Nil {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "slotId is required"))
		return
	}

	createLinkFlag := false
	createLinkFlag = body.CreateConferenceLink

	b, warning, err := h.UC.Create(r.Context(), au, body.SlotId, createLinkFlag)
	if err != nil {
		jsonresp.WriteAppError(w, err)
		return
	}
	if warning != "" {
		w.Header().Set("X-Booking-Warning", warning)
	}

	jsonresp.WriteJSON(w, http.StatusCreated, openapi.CreateBookingResponse{Booking: converter.BookingToOAPI(b)})
}
