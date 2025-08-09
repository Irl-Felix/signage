INSERT INTO Role (role_code, name, scope) VALUES
('SYS_ADMIN', 'System Administrator', 'global'),
('BUSINESS_OWNER', 'Business Owner', 'business'),
('BRANCH_MANAGER', 'Branch Manager', 'branch'),
('CONTENT_EDITOR', 'Content Editor', 'branch');

INSERT INTO Permission (code, description) VALUES
('MANAGE_USERS', 'Can manage user accounts'),
('VIEW_REPORTS', 'Can view analytics reports'),
('EDIT_CONTENT', 'Can create and edit content'),
('DELETE_CONTENT', 'Can delete content'),
('MANAGE_BILLING', 'Can manage subscription and billing');

INSERT INTO RolePermission (role_id, permission_id)
SELECT r.id, p.id FROM Role r, Permission p
WHERE
    (r.role_code='SYS_ADMIN' AND p.code IN ('MANAGE_USERS','VIEW_REPORTS','EDIT_CONTENT','DELETE_CONTENT','MANAGE_BILLING'))
    OR (r.role_code='BUSINESS_OWNER' AND p.code IN ('MANAGE_USERS','VIEW_REPORTS','EDIT_CONTENT','MANAGE_BILLING'))
    OR (r.role_code='BRANCH_MANAGER' AND p.code IN ('VIEW_REPORTS','EDIT_CONTENT','DELETE_CONTENT'))
    OR (r.role_code='CONTENT_EDITOR' AND p.code IN ('EDIT_CONTENT'));