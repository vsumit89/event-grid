package eventHandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/internal/commons"
	"server/internal/mocks"
	"server/internal/models"
	"server/pkg/logger"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
)

func getDummyEventDetails() *models.Event {
	userID := uint(1)
	eventID := uint(1)
	startTimestr := "2024-05-27 18:30:00"
	endTimeStr := "2024-05-27 19:30:00"

	start, _ := time.Parse("2006-01-02 15:04:05", startTimestr)

	end, _ := time.Parse("2006-01-02 15:04:05", endTimeStr)

	return &models.Event{
		Base: models.Base{
			ID: eventID,
		},
		CreatedBy:   userID,
		Title:       "test event",
		Description: "test description",
		Start:       start,
		End:         end,
		Attendees: []models.User{
			{
				Base: models.Base{
					ID: userID,
				},
			},
		},
	}
}

func TestGetEvent(t *testing.T) {
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
	// dummy data for testing

	successEvent := getDummyEventDetails()

	t.Run("TestGetEvent_Success", func(t *testing.T) {
		mockEventSvc.MockGetEventByID = func(userID, eventID uint) (*models.Event, error) {
			return successEvent, nil
		}

		parentCtx := context.Background()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("eventID", fmt.Sprintf("%d", successEvent.ID))

		ctx := context.WithValue(parentCtx, commons.ClaimsContext, &commons.CustomClaims{
			UserID: successEvent.CreatedBy,
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/events/{eventID}", nil)
		if err != nil {
			t.Error(err)
		}

		req = req.WithContext(ctx)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		eventHandler.getEvent(rr, req)

		if status := rr.Code; status != http.StatusOK {
			data := rr.Body.String()

			logger.Info("response", "data", data)

			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

	})

	t.Run("TestGetEvent_WithInvalidEventID", func(t *testing.T) {
		parentCtx := context.Background()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("eventID", "s")

		ctx := context.WithValue(parentCtx, commons.ClaimsContext, &commons.CustomClaims{
			UserID: successEvent.CreatedBy,
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/events/{eventID}", nil)
		if err != nil {
			t.Error(err)
		}

		req = req.WithContext(ctx)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		eventHandler.getEvent(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			data := rr.Body.String()

			logger.Info("response", "data", data)

			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("TestGetEvent_JWTClaimsFailed", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("eventID", fmt.Sprintf("%d", successEvent.ID))

		req, err := http.NewRequest("GET", "/api/v1/events/{eventID}", nil)
		if err != nil {
			t.Error(err)
		}

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		eventHandler.getEvent(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			data := rr.Body.String()

			logger.Info("response", "data", data)

			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("TestGetEvent_GetEventByIDFailed", func(t *testing.T) {
		mockEventSvc.MockGetEventByID = func(userID, eventID uint) (*models.Event, error) {
			return nil, errors.New("failed")
		}

		parentCtx := context.Background()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("eventID", fmt.Sprintf("%d", successEvent.ID))

		ctx := context.WithValue(parentCtx, commons.ClaimsContext, &commons.CustomClaims{
			UserID: successEvent.CreatedBy,
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/events/{eventID}", nil)
		if err != nil {
			t.Error(err)
		}

		req = req.WithContext(ctx)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		eventHandler.getEvent(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			data := rr.Body.String()

			logger.Info("response", "data", data)

			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)

		}
	})

	t.Run("TestGetEvent_EventNotFound", func(t *testing.T) {
		mockEventSvc.MockGetEventByID = func(userID, eventID uint) (*models.Event, error) {
			return nil, commons.ErrEventNotFound
		}

		parentCtx := context.Background()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("eventID", fmt.Sprintf("%d", successEvent.ID))

		ctx := context.WithValue(parentCtx, commons.ClaimsContext, &commons.CustomClaims{
			UserID: successEvent.CreatedBy,
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/events/{eventID}", nil)
		if err != nil {
			t.Error(err)
		}

		req = req.WithContext(ctx)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		eventHandler.getEvent(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			data := rr.Body.String()

			logger.Info("response", "data", data)

			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}
