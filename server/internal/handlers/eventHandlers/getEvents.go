package eventHandlers

import (
	"net/http"
	"server/internal/handlers/middlewares"
	"server/internal/services"
	"server/pkg/utils"
	"time"
)

func (h *Handler) getEvents(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetJWTClaims(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, err)
		return
	}

	startDate := r.URL.Query().Get("start_date")

	endDate := r.URL.Query().Get("end_date")

	var startTime, endTime time.Time

	if startDate == "" {
		startTime = time.Now().In(h.timezone)
	} else {
		startTime, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, err)
			return
		}
	}

	if endDate == "" {
		endTime = startTime.Add(24 * time.Hour)
	} else {
		endTime, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, err)
			return
		}

		// adding 24 hours to the end time
		endTime = endTime.Add(24 * time.Hour)
	}

	events, err := h.eventSvc.GetEvents(user.UserID, services.EventFilters{StartTime: startTime, EndTime: endTime})
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, events)
}
