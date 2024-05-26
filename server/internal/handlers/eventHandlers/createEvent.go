package eventHandlers

import (
	"encoding/json"
	"net/http"
	"server/internal/handlers/dtos"
	"server/internal/handlers/middlewares"
	"server/pkg/utils"
)

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetJWTClaims(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err)
		return
	}

	var req dtos.EventDTO

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = req.Validate()
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	event, err := h.eventSvc.CreateEvent(user.UserID, &req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	for i := 0; i < len(event.Attendees); i++ {
		event.Attendees[i].Password = ""
	}

	utils.SendJSON(w, http.StatusOK, event)
}
