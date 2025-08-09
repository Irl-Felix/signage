
-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Business (Tenants)
CREATE TABLE Business (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  subscription_plan TEXT NOT NULL DEFAULT 'basic'
    CHECK (subscription_plan IN ('basic', 'pro', 'enterprise')),
  status TEXT NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'active', 'suspended')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Location (Branch/Store)
CREATE TABLE Location (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id UUID NOT NULL REFERENCES Business(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  address TEXT,
  timezone TEXT
);