package userHandlers

import (
	"net/http"
	"server/internal/handlers/middlewares"
	"server/pkg/utils"
	"strconv"
)

func (h *Handler) searchUsers(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetJWTClaims(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err)
		return
	}

	query := r.URL.Query().Get("query")

	limit := r.URL.Query().Get("limit")

	var intLimit int

	if limit == "" {
		intLimit = 10
	} else {
		intLimit, err = strconv.Atoi(limit)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, err)
			return
		}
	}

	users, err := h.userSvc.GetUsers(user.UserID, query, intLimit)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	for i := 0; i < len(users); i++ {
		users[i].Password = ""
	}

	utils.SendJSON(w, http.StatusOK, users)
}
