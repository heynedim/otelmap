package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	errorz "github.com/jack5341/otel-map-server/internal/errors"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

func TestServiceMapGetInvalidToken(t *testing.T) {
	h := NewServiceMapHandler(nil, otel.Tracer("test"))
	e := echo.New()
	e.GET("/service-map/:session-token", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/service-map/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), errorz.ErrInvalidSessionToken.Error()) {
		t.Errorf("body = %s, want invalid token error", rec.Body.String())
	}
}

func TestServiceMapGetEmptyToken(t *testing.T) {
	h := NewServiceMapHandler(nil, otel.Tracer("test"))
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/service-map/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("session-token", "")

	if err := h.Get(c); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), errorz.ErrSessionTokenRequired.Error()) {
		t.Errorf("body = %s, want required token error", rec.Body.String())
	}
}
