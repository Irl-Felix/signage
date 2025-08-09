CREATE TABLE UserRoleAssignment (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES UserProfile(id) ON DELETE CASCADE,
  role_id UUID NOT NULL REFERENCES Role(id) ON DELETE CASCADE,
  business_id UUID REFERENCES Business(id),
  location_id UUID REFERENCES Location(id),
  is_owner BOOLEAN DEFAULT FALSE,
  assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (user_id, role_id, business_id, location_id)
);