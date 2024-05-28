package eventHandlers

import (
	"bytes"
	"context"
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

	"github.com/dgrijalva/jwt-go"
)

func TestCreateEvent(t *testing.T) {
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

	dummyEventDto := dtos.EventDTO{
		Title:       successEvent.Title,
		Description: successEvent.Description,
		StartTime:   successEvent.Start,
		EndTime:     successEvent.End,
		Attendees: []string{
			"test@gmailcom",
			"someemail@gmail.com",
		},
	}

	t.Run("TestCreateEvent_Sucess", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: successEvent.CreatedBy,
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})

		data, err := json.Marshal(&dummyEventDto)
		if err != nil {
			t.Error(err)
		}

		requestData := bytes.NewBuffer(data)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/events", requestData)
		if err != nil {
			t.Error(err)
		}

		mockEventSvc.MockCreateEvent = func(userID uint, req *dtos.EventDTO) (*models.Event, error) {
			return successEvent, nil
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.createEvent)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		t.Logf("handler returned response: %v", rr.Body.String())
	})

	t.Run("TestCreateEvent_ServiceError", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: successEvent.CreatedBy,
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})
		data, err := json.Marshal(&dummyEventDto)
		if err != nil {
			t.Error(err)
		}
		requestData := bytes.NewBuffer(data)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/events", requestData)
		if err != nil {
			t.Error(err)
		}

		mockEventSvc.MockCreateEvent = func(userID uint, req *dtos.EventDTO) (*models.Event, error) {
			return nil, errors.New("service error")
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.createEvent)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

	})

	t.Run("TestCreateEvent_Unauthorized", func(t *testing.T) {
		ctx := context.Background()
		data, err := json.Marshal(&dummyEventDto)
		if err != nil {
			t.Error(err)
		}
		requestData := bytes.NewBuffer(data)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/events", requestData)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.createEvent)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}

		// Add additional assertions for the response body, if needed
	})

	t.Run("TestCreateEvent_InvalidJSONBody", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: successEvent.CreatedBy,
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})

		requestData := bytes.NewBufferString("invalid json")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/events", requestData)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.createEvent)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

	})

	t.Run("TestCreateEvent_InvalidEventDTO", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), commons.ClaimsContext, &commons.CustomClaims{
			UserID: successEvent.CreatedBy,
			Email:  userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		})

		invalidEventDto := dtos.EventDTO{
			Title:       "",
			Description: successEvent.Description,
			StartTime:   successEvent.Start,
			EndTime:     successEvent.End,
			Attendees:   []string{"test@gmailcom", "someemail@gmail.com"},
		}
		data, err := json.Marshal(&invalidEventDto)
		if err != nil {
			t.Error(err)
		}
		requestData := bytes.NewBuffer(data)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/events", requestData)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventHandler.createEvent)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

	})
}
