package schedules

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/converter"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/middleware"
	schedulesuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/schedules"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

type Handler struct {
	UC *schedulesuc.UseCase
}

func NewHandler(uc *schedulesuc.UseCase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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
	var body openapi.Schedule
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "invalid json"))
		return
	}
	if body.RoomId != uuid.Nil && body.RoomId != roomID {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "roomId in body must match path"))
		return
	}
	sch, err := h.UC.Create(r.Context(), au, roomID, body.DaysOfWeek, body.StartTime, body.EndTime)
	if err != nil {
		jsonresp.WriteAppError(w, err)
		return
	}
	jsonresp.WriteJSON(w, http.StatusCreated, openapi.CreateScheduleResponse{Schedule: converter.ScheduleToOAPI(sch)})
}
