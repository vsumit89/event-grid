package userHandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/internal/commons"
	"server/internal/mocks"
	"server/internal/models"
	"server/pkg/logger"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestSearchUsers(t *testing.T) {
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

	t.Run("TestSearchUsers_Success", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/users/search?query=test&limit=5", nil)
		if err != nil {
			t.Error(err)
		}

		mockUserSvc.MockGetUsers = func(userID uint, query string, limit int) ([]models.User, error) {
			return []models.User{
				{
					Base: models.Base{
						ID: 1,
					},
					Name:  "John Doe",
					Email: "john@example.com",
				},
				{
					Base: models.Base{
						ID: 2,
					},
					Name:  "Jane Smith",
					Email: "jane@example.com",
				},
			}, nil
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.searchUsers)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("TestSearchUsers_Unauthorized", func(t *testing.T) {
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/users/search", nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.searchUsers)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("TestSearchUsers_InvalidLimit", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/users/search?query=test&limit=invalid", nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.searchUsers)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

	})

	t.Run("TestSearchUsers_ServiceError", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/users/search?query=test&limit=5", nil)
		if err != nil {
			t.Error(err)
		}

		mockUserSvc.MockGetUsers = func(userID uint, query string, limit int) ([]models.User, error) {
			return nil, errors.New("service error")
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.searchUsers)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

	})
}
