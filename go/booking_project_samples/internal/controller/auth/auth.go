package auth

import (
	"encoding/json"
	"net/http"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	authuc "github.com/test-backend-ddarvih/room-booking/internal/usercase/auth"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
	"github.com/test-backend-ddarvih/room-booking/pkg/generated/openapi"
)

func DummyLogin(authUC *authuc.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body openapi.DummyLoginJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonresp.WriteAppError(w, apperrors.BadRequest(apperrors.InvalidRequest, "invalid json"))
			return
		}

		if err := authuc.ValidateDummyRole(body.Role); err != nil {
			jsonresp.WriteAppError(w, err)
			return
		}

		uid, err := authuc.DummyUserIDForRole(body.Role)
		if err != nil {
			jsonresp.WriteAppError(w, apperrors.Internal("token issue failed"))
			return
		}

		token, err := authUC.GenerateAccessToken(uid, body.Role)
		if err != nil {
			jsonresp.WriteAppError(w, apperrors.Internal("token issue failed"))
			return
		}

		jsonresp.WriteJSON(w, http.StatusOK, openapi.Token{Token: token})
	}
}
