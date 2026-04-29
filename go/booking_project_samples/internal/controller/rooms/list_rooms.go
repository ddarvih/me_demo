package rooms

import (
	"net/http"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/converter"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/middleware"
	roomsuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/rooms"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

type Handler struct {
	UC *roomsuc.UseCase
}

func NewHandler(uc *roomsuc.UseCase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	au, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("unauthorized"))
		return
	}
	list, err := h.UC.List(r.Context(), au)
	if err != nil {
		jsonresp.WriteAppError(w, err)
		return
	}
	resp := openapi.ListRoomsResponse{Rooms: make([]openapi.Room, 0, len(list))}
	for _, room := range list {
		resp.Rooms = append(resp.Rooms, converter.RoomToOAPI(room))
	}
	jsonresp.WriteJSON(w, http.StatusOK, resp)
}
