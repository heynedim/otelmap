package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != "8000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8000")
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.ServiceName != "default" {
		t.Errorf("ServiceName = %q, want %q", cfg.ServiceName, "default")
	}
	if cfg.OTLPHTTPURL != "http://localhost/v1/traces" {
		t.Errorf("OTLPHTTPURL = %q, want default", cfg.OTLPHTTPURL)
	}
	if cfg.OTLPGRPCURL != "http://localhost:4317" {
		t.Errorf("OTLPGRPCURL = %q, want default", cfg.OTLPGRPCURL)
	}
	if cfg.CORSAllowedOrigins != "*" {
		t.Errorf("CORSAllowedOrigins = %q, want %q", cfg.CORSAllowedOrigins, "*")
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 10s", cfg.ShutdownTimeout)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "3")
	t.Setenv("OTLP_HTTP_URL", "https://otlp.example.com/v1/traces")
	t.Setenv("OTLP_GRPC_URL", "https://otlp.example.com:4317")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://a.com, https://b.com")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != "9000" {
		t.Errorf("Port = %q, want 9000", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", cfg.LogLevel)
	}
	if cfg.ShutdownTimeout != 3*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 3s", cfg.ShutdownTimeout)
	}
	if cfg.OTLPHTTPURL != "https://otlp.example.com/v1/traces" {
		t.Errorf("OTLPHTTPURL = %q", cfg.OTLPHTTPURL)
	}
	if cfg.OTLPGRPCURL != "https://otlp.example.com:4317" {
		t.Errorf("OTLPGRPCURL = %q", cfg.OTLPGRPCURL)
	}
	if cfg.CORSAllowedOrigins != "https://a.com, https://b.com" {
		t.Errorf("CORSAllowedOrigins = %q", cfg.CORSAllowedOrigins)
	}
}
