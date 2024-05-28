package eventHandlers

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

func TestGetEvents(t *testing.T) {
	logger.InitLogger()
	mockEventSvc := &mocks.MockEventSvc{}
	eventHandler := Handler{
		eventSvc: mockEventSvc,
		jwtSvc: &commons.JwtSvc{
			SecretKey: "my-secret",
			TTL:       time.Minute * 5,
		},
		timezone: time.FixedZone(commons.IST_TIMEZONE, commons.IST_OFFSET),
	}
	userEmail := "test@gmail.com"

	t.Run("TestGetEvents_Success", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/events?start_date=2023-05-28&end_date=2023-05-29", nil)
		if err != nil {
			t.Error(err)
		}

		startTime, _ := time.Parse("2006-01-02", "2023-05-28")
		endTime, _ := time.Parse("2006-01-02", "2023-05-29")
		endTime = endTime.Add(24 * time.Hour)

		mockEventSvc.MockGetEvents = func(userID uint, filters commons.EventFilters) ([]models.Event, error) {
			return []models.Event{
				{
					Base: models.Base{
						ID: 1,
					},
					Title:     "Test Event",
					CreatedBy: userID,
					Start:     startTime,
					End:       endTime,
				},
			}, nil
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.getEvents)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("TestGetEvents_Unauthorized", func(t *testing.T) {
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/events", nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.getEvents)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}

	})

	t.Run("TestGetEvents_InvalidStartDate", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/events?start_date=invalid-date", nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.getEvents)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("TestGetEvents_InvalidEndDate", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/events?start_date=2023-05-28&end_date=invalid-date", nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.getEvents)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

	})

	t.Run("TestGetEvents_ServiceError", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: uint(1),
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/events?start_date=2023-05-28&end_date=2023-05-29", nil)
		if err != nil {
			t.Error(err)
		}

		mockEventSvc.MockGetEvents = func(userID uint, filters commons.EventFilters) ([]models.Event, error) {
			return nil, errors.New("service error")
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.getEvents)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

	})
}
