package userHandlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/handlers/dtos"
	"server/internal/mocks"
	"server/internal/models"
	"server/pkg/logger"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	logger.InitLogger()
	mockUserSvc := &mocks.MockUserSvc{}
	authHandler := AuthHandler{
		userSvc: mockUserSvc,
	}

	t.Run("TestRegisterUser_Success", func(t *testing.T) {
		registerReq := dtos.CreateUserReqeust{
			Name:     "John Doe",
			Email:    "test@example.com",
			Password: "password",
		}
		reqBody, err := json.Marshal(&registerReq)
		if err != nil {
			t.Errorf("failed to marshal request body: %v", err)
		}

		mockUserSvc.MockRegisterUser = func(user *models.User) (*models.User, error) {
			return user, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.registerUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resp map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		if err != nil {
			t.Errorf("failed to unmarshal response body: %v", err)
		}

		expectedMessage := "user registered successfully"
		if resp["message"] != expectedMessage {
			t.Errorf("handler returned unexpected message: got %v want %v", resp["message"], expectedMessage)
		}
	})

	t.Run("TestRegisterUser_InvalidRequest", func(t *testing.T) {
		invalidReq := []byte(`{"email":"test@example.com"}`)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(invalidReq))
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.registerUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("TestRegisterUser_ServiceError", func(t *testing.T) {
		registerReq := dtos.CreateUserReqeust{
			Name:     "John Doe",
			Email:    "test@example.com",
			Password: "password",
		}
		reqBody, err := json.Marshal(&registerReq)
		if err != nil {
			t.Errorf("failed to marshal request body: %v", err)
		}

		mockUserSvc.MockRegisterUser = func(user *models.User) (*models.User, error) {
			return nil, errors.New("service error")
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.registerUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}
