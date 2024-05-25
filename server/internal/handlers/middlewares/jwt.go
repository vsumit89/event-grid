package middlewares

import (
	"context"
	"errors"
	"net/http"
	"server/internal/commons"
	"server/pkg/utils"
	"time"
)

func JWTAuth(svc *commons.JwtSvc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				utils.SendError(w, http.StatusUnauthorized, errors.New("you are unauthorized user, please login"))
				return
			}

			claims, err := svc.ValidateToken(cookie.Value)
			if err != nil {
				utils.SendError(w, http.StatusUnauthorized, commons.ErrTokenInvalid)
				return
			}

			if time.Now().Unix() > claims.ExpiresAt {
				utils.SendError(w, http.StatusUnauthorized, commons.ErrTokenExpired)
				return
			}

			ctx := context.WithValue(r.Context(), commons.ClaimsContext, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
