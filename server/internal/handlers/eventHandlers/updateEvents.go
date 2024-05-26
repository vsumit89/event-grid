package eventHandlers

import (
	"encoding/json"
	"net/http"
	"server/internal/handlers/dtos"
	"server/internal/handlers/middlewares"
	"server/pkg/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) updateEvent(w http.ResponseWriter, r *http.Request) {
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

	updatedEvent, err := h.eventSvc.UpdateEvent(user.UserID, uint(eventID), &req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, updatedEvent)
}
