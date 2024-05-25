package utils

import (
	"encoding/json"
	"net/http"
	"server/pkg/logger"
)

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error("error while sending response", "error", err.Error())
	}
}

type ErrResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

func SendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errBody := ErrResponse{
		Error{Message: err.Error()},
	}

	err = json.NewEncoder(w).Encode(errBody)
	if err != nil {
		logger.Error("error while sending response", "error", err.Error())
	}
}
