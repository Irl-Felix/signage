-- +migrate Up

-- Seed Roles
INSERT INTO Role (id, role_code, name, scope) VALUES
  (gen_random_uuid(), 'SYS_ADMIN', 'System Administrator', 'global'),
  (gen_random_uuid(), 'ACCT_MANAGER', 'Account Manager', 'global'),
  (gen_random_uuid(), 'TECH_SUPPORT', 'Technical Support Specialist', 'global'),
  (gen_random_uuid(), 'BILLING_ADMIN', 'Billing Administrator', 'global'),
  (gen_random_uuid(), 'ACCOUNT_ADMIN', 'Account Admin', 'business'),
  (gen_random_uuid(), 'CONTENT_MANAGER', 'Content Manager', 'business'),
  (gen_random_uuid(), 'LOCATION_MANAGER', 'Location Manager', 'business'),
  (gen_random_uuid(), 'LOCAL_UPLOADER', 'Local Content Uploader', 'business'),
  (gen_random_uuid(), 'VIEWER', 'Viewer', 'business');

-- Seed Permissions
INSERT INTO Permission (code, description) VALUES
  ('MANAGE_USERS', 'Manage all users'),
  ('VIEW_DASHBOARD', 'View dashboard'),
  ('EDIT_CONTENT', 'Edit content');

-- Map all permissions to SYS_ADMIN
INSERT INTO RolePermission (role_id, permission_id)
SELECT r.id, p.id FROM Role r, Permission p WHERE r.role_code = 'SYS_ADMIN';

-- Seed Business
INSERT INTO Business (id, name, subscription_plan)
VALUES ('11111111-1111-1111-1111-111111111111', 'Acme Corp', 'enterprise');

-- Seed Locations
INSERT INTO Location (id, business_id, name, address, timezone)
VALUES
  ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'Downtown HQ', '123 Main St, City', 'Asia/Taipei'),
  ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', 'Branch Office', '456 Market St, City', 'Asia/Taipei');

-- Seed Users
INSERT INTO UserProfile (id, supabase_uid, first_name, last_name, email, phone, profile_picture_url)
VALUES
  ('aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1', '61a10953-22b2-47ae-b041-4f73374ec140', 'Alice', 'Walker', 'alice.walker@acme.com', '1234567890', NULL),
  ('aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2', '72b4c6d1-77ef-4f28-92dd-101b27a899a3', 'Bob', 'Martinez', 'bob.martinez@acme.com', '2345678901', NULL),
  ('aaaaaaa3-aaaa-aaaa-aaaa-aaaaaaaaaaa3', 'f91e1a65-d74e-4cc1-84a7-91dbb3783bc5', 'Charlie', 'Nguyen', 'charlie.nguyen@acme.com', '3456789012', NULL);

-- Assign Global Role to Alice (System Administrator)
INSERT INTO UserRoleAssignment (id, user_id, role_id, assigned_at)
SELECT gen_random_uuid(), 'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1', r.id, NOW()
FROM Role r WHERE r.role_code = 'SYS_ADMIN';

-- Assign Business Role to Bob (Account Admin for Acme)
INSERT INTO UserRoleAssignment (id, user_id, role_id, business_id, assigned_at)
SELECT gen_random_uuid(), 'aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2', r.id, '11111111-1111-1111-1111-111111111111', NOW()
FROM Role r WHERE r.role_code = 'ACCOUNT_ADMIN';

-- Assign Location Role to Charlie (Viewer at Branch Office)
INSERT INTO UserRoleAssignment (id, user_id, role_id, business_id, location_id, assigned_at)
SELECT gen_random_uuid(), 'aaaaaaa3-aaaa-aaaa-aaaa-aaaaaaaaaaa3', r.id, '11111111-1111-1111-1111-111111111111', '33333333-3333-3333-3333-333333333333', NOW()
FROM Role r WHERE r.role_code = 'VIEWER';

-- Session log for Alice
INSERT INTO SessionLog (
  id, user_id, login_at, logout_at, ip_address, user_agent,
  access_token, refresh_token, is_active
) VALUES (
  gen_random_uuid(),
  'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
  NOW() - INTERVAL '1 hour',
  NOW(),
  '192.168.1.101',
  'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)',
  'access_token_example123',
  'refresh_token_example456',
  FALSE
);