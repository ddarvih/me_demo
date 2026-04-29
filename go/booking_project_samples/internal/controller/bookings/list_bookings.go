package bookings

import (
	"net/http"
	"strconv"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/converter"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	"github.com/test-backend-ddarvih/room-booking/internal/controller/middleware"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	au, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("unauthorized"))
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	list, pag, err := h.UC.ListAll(r.Context(), au, page, pageSize)
	if err != nil {
		jsonresp.WriteAppError(w, err)
		return
	}
	resp := openapi.ListBookingsResponse{
		Bookings:   make([]openapi.Booking, 0, len(list)),
		Pagination: converter.PaginationToOpenAPI(pag.Page, pag.PageSize, pag.Total),
	}
	for _, b := range list {
		resp.Bookings = append(resp.Bookings, converter.BookingToOAPI(b))
	}
	jsonresp.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) ListMy(w http.ResponseWriter, r *http.Request) {
	au, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("unauthorized"))
		return
	}
	list, err := h.UC.ListMy(r.Context(), au)
	if err != nil {
		jsonresp.WriteAppError(w, err)
		return
	}
	resp := openapi.ListMyBookingsResponse{Bookings: make([]openapi.Booking, 0, len(list))}
	for _, b := range list {
		resp.Bookings = append(resp.Bookings, converter.BookingToOAPI(b))
	}
	jsonresp.WriteJSON(w, http.StatusOK, resp)
}
