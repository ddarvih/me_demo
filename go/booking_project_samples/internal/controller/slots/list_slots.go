package slots

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/converter"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/middleware"
	slotsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/slots"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

type Handler struct {
	UC *slotsuc.UseCase
}

func NewHandler(uc *slotsuc.UseCase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	au, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("unauthorized"))
		return
	}
	roomID, err := uuid.Parse(chi.URLParam(r, "roomId"))
	if err != nil {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "invalid room id"))
		return
	}
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "date query parameter is required"))
		return
	}
	dateUTC, err := time.ParseInLocation("2006-01-02", dateStr, time.UTC)
	if err != nil {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "invalid date format, use YYYY-MM-DD"))
		return
	}
	slots, err := h.UC.ListAvailable(r.Context(), au, roomID, dateUTC)
	if err != nil {
		jsonresp.WriteAppError(w, err)
		return
	}
	resp := openapi.ListSlotsResponse{Slots: make([]openapi.Slot, 0, len(slots))}
	for _, s := range slots {
		resp.Slots = append(resp.Slots, converter.SlotToOAPI(s))
	}
	jsonresp.WriteJSON(w, http.StatusOK, resp)
}
