package jsonresp

import (
	"encoding/json"
	"net/http"

	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteAppError(w http.ResponseWriter, err error) {
	ae := apperrors.AsAppError(err)
	if ae == nil {
		return
	}
	body := openapi.ErrorResponse{
		Error: openapi.ErrorResponseError{
			Code:    string(ae.Code),
			Message: ae.Message,
		},
	}
	WriteJSON(w, ae.Status, body)
}

func WriteInternal(w http.ResponseWriter) {
	body := openapi.InternalErrorResponse{
		Error: openapi.InternalErrorResponseError{
			Code:    string(apperrors.InternalError),
			Message: "internal server error",
		},
	}
	WriteJSON(w, http.StatusInternalServerError, body)
}
