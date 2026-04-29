package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/test-backend-ddarvih/room-booking/internal/controller/jsonresp"
	"github.com/test-backend-ddarvih/room-booking/internal/entity"
	"github.com/test-backend-ddarvih/room-booking/pkg/apperrors"
)

type ctxKey string

const authUserKey ctxKey = "authUser"

func OnlyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		au, ok := r.Context().Value(authUserKey).(entity.AuthUser)
		if !ok || au.Role != entity.RoleAdmin {
			jsonresp.WriteAppError(w, apperrors.ForbiddenMsg("admin role required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func OnlyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		au, ok := r.Context().Value(authUserKey).(entity.AuthUser)
		if !ok || au.Role != entity.RoleUser {
			jsonresp.WriteAppError(w, apperrors.ForbiddenMsg("user role required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func BearerJWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") || parts[1] == "" {
				jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("missing bearer token"))
				return
			}
			raw := strings.TrimSpace(parts[1])
			tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, apperrors.UnauthorizedMsg("invalid signing method")
				}
				return []byte(secret), nil
			})
			if err != nil || !tok.Valid {
				jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("invalid token"))
				return
			}
			claims, ok := tok.Claims.(jwt.MapClaims)
			if !ok {
				jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("invalid claims"))
				return
			}
			uidStr, _ := claims["user_id"].(string)
			role, _ := claims["role"].(string)
			uid, err := uuid.Parse(uidStr)
			if err != nil || uid == uuid.Nil {
				jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("invalid user_id in token"))
				return
			}
			if role != entity.RoleAdmin && role != entity.RoleUser {
				jsonresp.WriteAppError(w, apperrors.UnauthorizedMsg("invalid role in token"))
				return
			}
			au := entity.AuthUser{ID: uid, Role: role}
			ctx := context.WithValue(r.Context(), authUserKey, au)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthUserFromContext(ctx context.Context) (entity.AuthUser, bool) {
	v, ok := ctx.Value(authUserKey).(entity.AuthUser)
	return v, ok
}
