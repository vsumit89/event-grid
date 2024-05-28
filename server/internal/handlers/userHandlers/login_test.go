package userHandlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/commons"
	"server/internal/handlers/dtos"
	"server/internal/mocks"
	"server/internal/models"
	"server/pkg/logger"
	"testing"
	"time"
)

func TestLoginUser(t *testing.T) {
	logger.InitLogger()
	mockUserSvc := &mocks.MockUserSvc{}
	authHandler := AuthHandler{
		userSvc: mockUserSvc,
		jwtSvc: &commons.JwtSvc{
			SecretKey: "my-secret",
			TTL:       time.Minute * 5,
		},
	}

	t.Run("TestLoginUser_Success", func(t *testing.T) {
		loginReq := dtos.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}
		reqBody, err := json.Marshal(&loginReq)
		if err != nil {
			t.Errorf("failed to marshal request body: %v", err)
		}

		expectedUser := &models.User{
			Base: models.Base{
				ID: 1,
			},
			Name:  "John Doe",
			Email: "test@example.com",
		}

		mockUserSvc.MockLogin = func(email, password string) (*models.User, error) {
			return expectedUser, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.loginUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		cookies := rr.Result().Cookies()
		if len(cookies) != 1 {
			t.Errorf("handler set incorrect number of cookies: got %d want 1", len(cookies))
		}

		var resp map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		if err != nil {
			t.Errorf("failed to unmarshal response body: %v", err)
		}

		expectedMessage := "user logged in successfully"
		if resp["message"] != expectedMessage {
			t.Errorf("handler returned unexpected message: got %v want %v", resp["message"], expectedMessage)
		}
	})

	t.Run("TestLoginUser_InvalidRequest", func(t *testing.T) {
		invalidReq := []byte(`{"email":"test@example.com"}`)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(invalidReq))
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.loginUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("TestLoginUser_ServiceError", func(t *testing.T) {
		loginReq := dtos.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}
		reqBody, err := json.Marshal(&loginReq)
		if err != nil {
			t.Errorf("failed to marshal request body: %v", err)
		}

		mockUserSvc.MockLogin = func(email, password string) (*models.User, error) {
			return nil, errors.New("service error")
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(authHandler.loginUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}
