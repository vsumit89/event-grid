package middlewares

import (
	"net/http"
	"server/internal/commons"
	"server/pkg/utils"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				utils.SendError(w, http.StatusInternalServerError, commons.ErrInternalServer)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
