-- +migrate Up
CREATE TABLE Role (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  role_code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  scope TEXT NOT NULL CHECK (scope IN ('global', 'business', 'branch'))
);

CREATE TABLE Permission (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  description TEXT
);

CREATE TABLE RolePermission (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  role_id UUID NOT NULL REFERENCES Role(id) ON DELETE CASCADE,
  permission_id UUID NOT NULL REFERENCES Permission(id) ON DELETE CASCADE,
  UNIQUE (role_id, permission_id)
);