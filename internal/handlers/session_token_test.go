package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jack5341/otel-map-server/internal/config"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

func TestSessionTokenCreate(t *testing.T) {
	db := newTestDB(t)
	cfg := &config.Config{
		OTLPHTTPURL: "http://localhost/v1/traces",
		OTLPGRPCURL: "http://localhost:4317",
	}
	h := NewSessionTokenHandler(db, otel.Tracer("test"), cfg)

	e := echo.New()
	e.POST("/session-token", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/session-token", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	var resp SessionTokenResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, err := uuid.Parse(resp.Token); err != nil {
		t.Errorf("token %q is not a valid UUID: %v", resp.Token, err)
	}
	if resp.Ingest.OTLPHTTPURL != cfg.OTLPHTTPURL {
		t.Errorf("OTLPHTTPURL = %q, want %q", resp.Ingest.OTLPHTTPURL, cfg.OTLPHTTPURL)
	}
	if resp.Ingest.OTLPGRPCURL != cfg.OTLPGRPCURL {
		t.Errorf("OTLPGRPCURL = %q, want %q", resp.Ingest.OTLPGRPCURL, cfg.OTLPGRPCURL)
	}
	if resp.Ingest.HeaderKey != "X-OTEL-SESSION" {
		t.Errorf("HeaderKey = %q, want X-OTEL-SESSION", resp.Ingest.HeaderKey)
	}
	if resp.Ingest.HeaderValue != resp.Token {
		t.Errorf("HeaderValue = %q, want token %q", resp.Ingest.HeaderValue, resp.Token)
	}
	if resp.Ingest.ResourceAttribute.Key != "otelmap.session_token" {
		t.Errorf("ResourceAttribute.Key = %q, want otelmap.session_token", resp.Ingest.ResourceAttribute.Key)
	}
	if resp.Ingest.ResourceAttribute.Value != resp.Token {
		t.Errorf("ResourceAttribute.Value = %q, want token %q", resp.Ingest.ResourceAttribute.Value, resp.Token)
	}
}
