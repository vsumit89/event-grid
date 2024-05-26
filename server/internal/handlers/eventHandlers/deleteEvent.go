package eventHandlers

import (
	"net/http"
	"server/internal/handlers/middlewares"
	"server/pkg/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetJWTClaims(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err)
		return
	}

	eventID, err := strconv.Atoi(chi.URLParam(r, "eventID"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = h.eventSvc.DeleteEvent(user.UserID, uint(eventID))
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Event deleted successfully",
	})
}
