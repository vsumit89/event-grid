package userHandlers

import (
	"encoding/json"
	"net/http"
	"server/internal/handlers/dtos"
	"server/pkg/utils"
)

func (h *AuthHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateUserReqeust

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

	user := req.MapToUserModel()

	_, err = h.userSvc.RegisterUser(user)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"message": "user registered successfully",
	})
}
