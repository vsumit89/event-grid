package eventHandlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/internal/commons"
	"server/internal/mocks"
	"server/pkg/logger"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
)

func TestDeleteEvent(t *testing.T) {
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

	successEvent := getDummyEventDetails()

	t.Run("Success", func(t *testing.T) {

		mockEventSvc.MockDeleteEvent = func(userID, eventID uint) error {
			return nil
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

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/api/v1/events/{eventID}", nil)
		if err != nil {
			t.Error(err)
		}

		req = req.WithContext(ctx)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		eventHandler.deleteEvent(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

	})

	t.Run("TestDeleteFailure", func(t *testing.T) {
		mockEventSvc.MockDeleteEvent = func(userID, eventID uint) error {
			return commons.ErrEventNotFound
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

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/api/v1/events/{eventID}", nil)
		if err != nil {
			t.Error(err)
		}

		req = req.WithContext(ctx)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		eventHandler.deleteEvent(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("TestDeleteEvent_JWTClaimsFailed", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("eventID", fmt.Sprintf("%d", successEvent.ID))

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/events/{eventID}", nil)
		if err != nil {
			t.Error(err)
		}

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rr := httptest.NewRecorder()

		eventHandler.deleteEvent(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			data := rr.Body.String()

			logger.Info("response", "data", data)

			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}

		t.Run("TestDeleteEvent_WithInvalidEventID", func(t *testing.T) {
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

			req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/api/v1/events/{eventID}", nil)
			if err != nil {
				t.Error(err)
			}

			req = req.WithContext(ctx)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			eventHandler.deleteEvent(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				data := rr.Body.String()

				logger.Info("response", "data", data)

				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}
		})
	})

}
