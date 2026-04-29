package rooms

import (
	"encoding/json"
	"net/http"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/converter"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/middleware"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	au, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("unauthorized"))
		return
	}
	var body openapi.CreateRoomJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "invalid json"))
		return
	}
	room, err := h.UC.Create(r.Context(), au, body.Name, body.Description, body.Capacity)
	if err != nil {
		jsonresp.WriteAppError(w, err)
		return
	}
	jsonresp.WriteJSON(w, http.StatusCreated, openapi.CreateRoomResponse{Room: converter.RoomToOAPI(room)})
}
