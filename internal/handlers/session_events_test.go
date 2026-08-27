package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jack5341/otel-map-server/internal/config"
	errorz "github.com/jack5341/otel-map-server/internal/errors"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

func TestSessionEventsMissingToken(t *testing.T) {
	h := NewSessionEventsHandler(nil, otel.Tracer("test"), &config.Config{})
	e := echo.New()
	e.GET("/session-events", h.Listen)

	req := httptest.NewRequest(http.MethodGet, "/session-events", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), errorz.ErrSessionTokenRequired.Error()) {
		t.Errorf("body = %s, want required token error", rec.Body.String())
	}
}

func TestSessionEventsInvalidToken(t *testing.T) {
	h := NewSessionEventsHandler(nil, otel.Tracer("test"), &config.Config{})
	e := echo.New()
	e.GET("/session-events", h.Listen)

	req := httptest.NewRequest(http.MethodGet, "/session-events?token=not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), errorz.ErrInvalidSessionToken.Error()) {
		t.Errorf("body = %s, want invalid token error", rec.Body.String())
	}
}
