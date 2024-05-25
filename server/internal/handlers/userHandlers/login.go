package userHandlers

import (
	"encoding/json"
	"net/http"
	"server/internal/handlers/dtos"
	"server/pkg/logger"
	"server/pkg/utils"
	"time"
)

func (h *AuthHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	var req dtos.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("failed to generate token", "error", err)
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = req.Validate()
	if err != nil {
		logger.Error("failed to generate token", "error", err)
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.userSvc.Login(req.Email, req.Password)
	if err != nil {
		logger.Error("failed to generate token", "error", err)
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	token, err := h.jwtSvc.GenerateToken(user.Email, user.ID)
	if err != nil {
		logger.Error("failed to generate token", "error", err)
		utils.SendError(w, http.StatusInternalServerError, err)
		return
	}

	// expiresAfter is the time after which the token expires
	// it's in time.Duration
	expiresAfter := h.jwtSvc.GetTTL()

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(expiresAfter),
		HttpOnly: true,
	})

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"message": "user logged in successfully",
	})
}
