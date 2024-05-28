package userHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"server/internal/commons"
	"server/internal/mocks"
	"server/internal/models"
	"server/pkg/logger"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGetUser(t *testing.T) {
	logger.InitLogger()
	mockUserSvc := &mocks.MockUserSvc{}
	userHandler := Handler{
		userSvc: mockUserSvc,
		jwtSvc: &commons.JwtSvc{
			SecretKey: "my-secret",
			TTL:       time.Minute * 5,
		},
	}
	userEmail := "test@gmail.com"

	t.Run("TestGetUser_Success", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/users/me", nil)
		if err != nil {
			t.Error(err)
		}

		expectedUser := &models.User{
			Base: models.Base{
				ID: 1,
			},
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockUserSvc.MockGetUserByID = func(userID uint) (*models.User, error) {
			return expectedUser, nil
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.getUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var actualUser *models.User
		err = json.Unmarshal(rr.Body.Bytes(), &actualUser)
		if err != nil {
			t.Errorf("failed to unmarshal response body: %v", err)
		}

		actualUser.Password = "" // Remove password field for comparison
		if !reflect.DeepEqual(actualUser, expectedUser) {
			t.Errorf("handler returned unexpected user: got %v want %v", actualUser, expectedUser)
		}
	})

	t.Run("TestGetUser_Unauthorized", func(t *testing.T) {
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/users/me", nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.getUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("TestGetUser_ServiceError", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/users/me", nil)
		if err != nil {
			t.Error(err)
		}

		mockUserSvc.MockGetUserByID = func(userID uint) (*models.User, error) {
			return nil, errors.New("service error")
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.getUser)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}
