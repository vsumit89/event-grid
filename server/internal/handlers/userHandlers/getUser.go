package userHandlers

import (
	"net/http"
	"server/internal/handlers/middlewares"
	"server/pkg/utils"
)

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetJWTClaims(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err)
		return
	}

	userProfile, err := h.userSvc.GetUserByID(user.UserID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	userProfile.Password = ""

	utils.SendJSON(w, http.StatusOK, userProfile)
}
