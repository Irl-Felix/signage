-- +migrate Up

-- Optional: Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Business (Tenants)
CREATE TABLE Business (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  subscription_plan TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Location (Branch/Store)
CREATE TABLE Location (
  id UUID PRIMARY KEY,
  business_id UUID NOT NULL REFERENCES Business(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  address TEXT,
  timezone TEXT
);

-- 3. UserProfile
CREATE TABLE UserProfile (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  supabase_uid UUID NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT,
  phone TEXT,
  profile_picture_url TEXT,
  preferred_language TEXT DEFAULT 'en',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Role (RBAC)
CREATE TABLE Role (
  id UUID PRIMARY KEY,
  role_code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  scope TEXT NOT NULL CHECK (scope IN ('global', 'business'))
);

-- 5. Permission
CREATE TABLE Permission (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  description TEXT
);

-- 6. RolePermission (Role <-> Permission mapping)
CREATE TABLE RolePermission (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  role_id UUID NOT NULL REFERENCES Role(id) ON DELETE CASCADE,
  permission_id UUID NOT NULL REFERENCES Permission(id) ON DELETE CASCADE
);

-- 7. UserRoleAssignment
CREATE TABLE UserRoleAssignment (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES UserProfile(id) ON DELETE CASCADE,
  role_id UUID NOT NULL REFERENCES Role(id) ON DELETE CASCADE,
  business_id UUID REFERENCES Business(id),
  location_id UUID REFERENCES Location(id),
  assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. SessionLog
CREATE TABLE SessionLog (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES UserProfile(id) ON DELETE CASCADE,
  login_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  logout_at TIMESTAMP,
  ip_address TEXT,
  user_agent TEXT,
  access_token TEXT,
  refresh_token TEXT,
  token_expires_at TIMESTAMP, -- when the token expires
  is_active BOOLEAN DEFAULT TRUE,
  duration_seconds BIGINT GENERATED ALWAYS AS (
    CASE
      WHEN logout_at IS NOT NULL THEN EXTRACT(EPOCH FROM (logout_at - login_at))::BIGINT
      ELSE NULL
    END
  ) STORED
);