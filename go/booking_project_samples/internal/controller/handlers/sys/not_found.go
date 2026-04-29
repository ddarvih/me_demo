package sys

import (
	"io"
	"net/http"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
)

func NotFoundOrStub(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && (r.URL.Path == "/register" || r.URL.Path == "/login") {
		jsonresp.WriteJSON(w, http.StatusNotFound, map[string]string{
			"error": "not found",
		})
		return
	}
	_, _ = io.Copy(io.Discard, r.Body)
	_ = r.Body.Close()
	w.WriteHeader(http.StatusNotFound)
}
