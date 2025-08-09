-- +migrate Up
CREATE TABLE UserProfile (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  supabase_uid UUID NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT,
  phone TEXT,
  profile_picture_url TEXT,
  preferred_language TEXT DEFAULT 'en',
  status TEXT NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'active', 'suspended')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);