package middlewares

import (
	"context"
	"net/http"
	"server/internal/commons"
	"server/pkg/utils"
	"strings"
)

// parse the Bearer token from the header
func parseToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")

	if token == "" {
		return "", commons.ErrTokenNotFound
	}

	tokenDetails := strings.Split(token, " ")

	if len(tokenDetails) < 2 || tokenDetails[0] != "Bearer" {
		return "", commons.ErrTokenNotFound
	}

	return tokenDetails[1], nil
}

func JWTAuth(svc *commons.JwtSvc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := parseToken(r)
			if err != nil {
				utils.SendError(w, http.StatusUnauthorized, err)
				return
			}

			claims, err := svc.ValidateToken(token)
			if err != nil {
				utils.SendError(w, http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), commons.ClaimsContext, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
