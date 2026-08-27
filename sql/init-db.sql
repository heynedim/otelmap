-- Session tokens live in the `default` database (the one the server DSN points at).
-- Created here with a TTL so stale sessions are purged automatically; the server's
-- GORM AutoMigrate is a no-op (IF NOT EXISTS) when this table already exists.
CREATE TABLE IF NOT EXISTS default.session_tokens
(
    `token`      UUID,
    `created_at` DateTime DEFAULT now(),
    `updated_at` DateTime DEFAULT now()
)
ENGINE = MergeTree
ORDER BY token
TTL created_at + INTERVAL 7 DAY;

-- NOTE: default.otel_traces is created automatically by the OTel Collector's
-- ClickHouse exporter (create_schema defaults to true) on first ingest.
