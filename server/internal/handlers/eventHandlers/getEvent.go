package eventHandlers

import (
	"errors"
	"net/http"
	"server/internal/commons"
	"server/internal/handlers/middlewares"
	"server/pkg/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) getEvent(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetJWTClaims(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err)
		return
	}

	eventID, err := strconv.Atoi(chi.URLParam(r, "eventID"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, commons.ErrInvalidEventID)
		return
	}

	event, err := h.eventSvc.GetEventByID(user.UserID, uint(eventID))
	if err != nil {
		if errors.Is(err, commons.ErrEventNotFound) {
			utils.SendError(w, http.StatusNotFound, err)
			return
		}
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	for i := 0; i < len(event.Attendees); i++ {
		event.Attendees[i].Password = ""
	}

	utils.SendJSON(w, http.StatusOK, event)
}
