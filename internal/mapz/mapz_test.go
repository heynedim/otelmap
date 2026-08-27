package mapz

import (
	"context"
	"errors"
	"testing"

	errorz "github.com/jack5341/otel-map-server/internal/errors"
	"go.opentelemetry.io/otel"
)

func TestGetEdgesRequiresToken(t *testing.T) {
	m := NewMapper(nil, otel.Tracer("test"), context.Background())
	if _, err := m.GetEdges(""); !errors.Is(err, errorz.ErrSessionTokenRequired) {
		t.Fatalf("GetEdges(\"\") error = %v, want ErrSessionTokenRequired", err)
	}
}

func TestGetServicesWithMetricsRequiresToken(t *testing.T) {
	m := NewMapper(nil, otel.Tracer("test"), context.Background())
	if _, err := m.GetServicesWithMetrics(""); !errors.Is(err, errorz.ErrSessionTokenRequired) {
		t.Fatalf("GetServicesWithMetrics(\"\") error = %v, want ErrSessionTokenRequired", err)
	}
}
