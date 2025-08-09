-- +migrate Up
CREATE TABLE SessionLog (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES UserProfile(id) ON DELETE CASCADE,
  login_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  logout_at TIMESTAMP,
  ip_address TEXT,
  user_agent TEXT,
  access_token TEXT,
  refresh_token TEXT,
  token_expires_at TIMESTAMP,
  is_active BOOLEAN DEFAULT TRUE,
  duration_seconds BIGINT GENERATED ALWAYS AS (
    CASE
      WHEN logout_at IS NOT NULL THEN EXTRACT(EPOCH FROM (logout_at - login_at))::BIGINT
      ELSE NULL
    END
  ) STORED
);