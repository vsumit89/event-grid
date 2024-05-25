package eventHandlers

import (
	"encoding/json"
	"net/http"
	"server/internal/handlers/dtos"
	"server/pkg/utils"
)

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateEvent

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = req.Validate()
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	eventResponse, err := h.eventSvc.CreateEvent(&req)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, eventResponse)
}
