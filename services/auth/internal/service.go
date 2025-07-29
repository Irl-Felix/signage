package internal

import (
	"database/sql"
	"time"

	"github.com/Irl-Felix/signage/shared/auth"
)

type AuthService struct {
	DB *sql.DB
}

// =============================
// Role Management
// =============================

func (s *AuthService) CreateRole(role *Role) error {
	return s.DB.QueryRow(`
		INSERT INTO Role (id, role_code, name, scope)
		VALUES (gen_random_uuid(), $1, $2, $3)
		RETURNING id`,
		role.RoleCode, role.Name, role.Scope,
	).Scan(&role.ID)
}

func (s *AuthService) ListRoles() ([]Role, error) {
	rows, err := s.DB.Query(`SELECT id, role_code, name, scope FROM Role`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.RoleCode, &r.Name, &r.Scope); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func (s *AuthService) GetRole(id string) (*Role, error) {
	var role Role
	err := s.DB.QueryRow(`SELECT id, role_code, name, scope FROM Role WHERE id = $1`, id).
		Scan(&role.ID, &role.RoleCode, &role.Name, &role.Scope)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *AuthService) UpdateRole(id, roleCode, name, scope string) (bool, error) {
	res, err := s.DB.Exec(`UPDATE Role SET role_code = $1, name = $2, scope = $3 WHERE id = $4`, roleCode, name, scope, id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

func (s *AuthService) DeleteRole(id string) (bool, error) {
	res, err := s.DB.Exec(`DELETE FROM Role WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// =============================
// Permission Management
// =============================

func (s *AuthService) AddRolePermission(roleCode, permissionCode string) error {
	_, err := s.DB.Exec(`
		INSERT INTO RolePermission (id, role_id, permission_id)
		SELECT gen_random_uuid(), r.id, p.id
		FROM Role r, Permission p
		WHERE r.role_code = $1 AND p.code = $2`,
		roleCode, permissionCode)
	return err
}

func (s *AuthService) RemovePermissionFromRole(roleID, permissionID string) (bool, error) {
	res, err := s.DB.Exec(`DELETE FROM RolePermission WHERE role_id = $1 AND permission_id = $2`, roleID, permissionID)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

func (s *AuthService) AssignPermissionsToRole(roleID string, permissionIDs []string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.Exec(`DELETE FROM RolePermission WHERE role_id = $1`, roleID)
	if err != nil {
		return err
	}

	for _, pid := range permissionIDs {
		if _, err = tx.Exec(`INSERT INTO RolePermission (role_id, permission_id) VALUES ($1, $2)`, roleID, pid); err != nil {
			return err
		}
	}
	return nil
}

func (s *AuthService) GetRolePermissions(roleID string) ([]auth.Permission, error) {
	rows, err := s.DB.Query(`
			   SELECT p.id, p.code
			   FROM Permission p
			   JOIN RolePermission rp ON p.id = rp.permission_id
			   WHERE rp.role_id = $1`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []auth.Permission
	for rows.Next() {
		var p auth.Permission
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (s *AuthService) ListRolePermissions() ([]struct{ RoleCode, PermissionCode string }, error) {
	rows, err := s.DB.Query(`
		SELECT r.role_code, p.code
		FROM RolePermission rp
		JOIN Role r ON rp.role_id = r.id
		JOIN Permission p ON rp.permission_id = p.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []struct{ RoleCode, PermissionCode string }
	for rows.Next() {
		var entry struct{ RoleCode, PermissionCode string }
		if err := rows.Scan(&entry.RoleCode, &entry.PermissionCode); err != nil {
			return nil, err
		}
		mappings = append(mappings, entry)
	}
	return mappings, nil
}

// =============================
// User & Role Assignment
// =============================

func (s *AuthService) AddUserRole(userID, roleCode string) error {
	// Only add the role if the user does not already have it
	var exists bool
	err := s.DB.QueryRow(`
			   SELECT EXISTS (
					   SELECT 1 FROM UserRoleAssignment ura
					   JOIN Role r ON ura.role_id = r.id
					   WHERE ura.user_id = $1 AND r.role_code = $2
			   )
	   `, userID, roleCode).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // Already assigned, do nothing
	}
	_, err = s.DB.Exec(`
			   INSERT INTO UserRoleAssignment (id, user_id, role_id, assigned_at)
			   SELECT gen_random_uuid(), $1, r.id, NOW()
			   FROM Role r WHERE r.role_code = $2`,
		userID, roleCode)
	return err
}

type UserRoleAssignment struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	RoleCode string `json:"role_code"`
}

func (s *AuthService) ListUserRoles() ([]UserRoleAssignment, error) {
	rows, err := s.DB.Query(`
	   SELECT ura.user_id, u.email, r.role_code
	   FROM UserRoleAssignment ura
	   JOIN Role r ON ura.role_id = r.id
	   JOIN UserProfile u ON ura.user_id = u.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []UserRoleAssignment
	for rows.Next() {
		var entry UserRoleAssignment
		if err := rows.Scan(&entry.UserID, &entry.Email, &entry.RoleCode); err != nil {
			return nil, err
		}
		assignments = append(assignments, entry)
	}
	return assignments, nil
}

// =============================
// User Profile Management
// =============================

func (s *AuthService) CreateUser(u *UserProfile) error {
	return s.DB.QueryRow(`
			   INSERT INTO UserProfile (
					   supabase_uid, first_name, last_name, email, phone, profile_picture_url, preferred_language
			   )
			   VALUES ($1, $2, $3, $4, $5, $6, $7)
			   RETURNING id`,
		u.SupabaseUID, u.FirstName, u.LastName, u.Email, u.Phone, u.ProfilePictureURL, u.PreferredLanguage,
	).Scan(&u.ID)
}

func (s *AuthService) UpdateUserProfile(u *UserProfile) error {
	_, err := s.DB.Exec(`
			   UPDATE UserProfile
			   SET first_name = $1, last_name = $2, email = $3, phone = $4, profile_picture_url = $5, preferred_language = $6
			   WHERE id = $7`,
		u.FirstName, u.LastName, u.Email, u.Phone, u.ProfilePictureURL, u.PreferredLanguage, u.ID)
	return err
}

func (s *AuthService) GetUserBySupabaseUID(uid string) (*UserProfile, error) {
	row := s.DB.QueryRow(`SELECT id, email FROM UserProfile WHERE supabase_uid = $1`, uid)
	var user UserProfile
	if err := row.Scan(&user.ID, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) ListUsers() ([]auth.UserProfilePlain, error) {
	rows, err := s.DB.Query(`
			   SELECT id, supabase_uid, first_name, last_name, email, phone, profile_picture_url, preferred_language, created_at
			   FROM UserProfile`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []auth.UserProfilePlain
	for rows.Next() {
		var u auth.UserProfile
		if err := rows.Scan(&u.ID, &u.SupabaseUID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.ProfilePictureURL, &u.PreferredLanguage, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u.ToPlain())
	}
	return users, nil
}

// =============================
// Session Management
// =============================

func (s *AuthService) CreateSession(log *SessionLog) error {
	_, err := s.DB.Exec(`
			   UPDATE SessionLog
			   SET is_active = false, logout_at = NOW()
			   WHERE user_id = $1 AND is_active = true`, log.UserID)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`
			   INSERT INTO SessionLog (
					   user_id, ip_address, user_agent, access_token, refresh_token, token_expires_at, is_active
			   )
			   VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		log.UserID, log.IPAddress, log.UserAgent, log.AccessToken, log.RefreshToken, log.TokenExpiresAt, log.IsActive)
	return err
}

func (s *AuthService) LogoutSession(accessToken string, logoutAt time.Time) error {
	_, err := s.DB.Exec(`
			   UPDATE SessionLog
			   SET logout_at = $1, is_active = false
			   WHERE access_token = $2 AND is_active = true`,
		logoutAt, accessToken)
	return err
}

// =============================
// Audit Logs
// =============================

type AuditLog struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	Details   string    `json:"details"`
}

func (s *AuthService) ListAuditLogs() ([]AuditLog, error) {
	rows, err := s.DB.Query(`
			   SELECT id, user_id, action, created_at, details
			   FROM AuditLog
			   ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var l AuditLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.CreatedAt, &l.Details); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

// =============================
// User Management (Single User)
// =============================

// DeleteUser deletes a user by their ID
func (s *AuthService) DeleteUser(id string) (bool, error) {
	res, err := s.DB.Exec(`DELETE FROM UserProfile WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// Permission struct for user permissions
type Permission struct {
	ID          string
	Code        string
	Description string
}

// GetUserPermissions returns all permissions for a user by user ID
func (s *AuthService) GetUserPermissions(userID string) ([]Permission, error) {
	rows, err := s.DB.Query(`
			   SELECT p.id, p.code, p.description
			   FROM UserRoleAssignment ura
			   JOIN RolePermission rp ON ura.role_id = rp.role_id
			   JOIN Permission p ON rp.permission_id = p.id
			   WHERE ura.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Description); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

type UserWithPermissions struct {
	ID                string       `json:"id"`
	Email             string       `json:"email"`
	FirstName         string       `json:"first_name"`
	LastName          string       `json:"last_name"`
	Phone             string       `json:"phone"`
	ProfilePictureURL string       `json:"profile_picture_url"`
	PreferredLanguage string       `json:"preferred_language"`
	Permissions       []Permission `json:"permissions"`
}

// ListUsersWithPermissions returns all users and their permissions
func (s *AuthService) ListUsersWithPermissions() ([]UserWithPermissions, error) {
	rows, err := s.DB.Query(`
		SELECT u.id, u.email, u.first_name, u.last_name, u.phone, u.profile_picture_url, u.preferred_language,
			   p.id, p.code, p.description
		FROM UserProfile u
		LEFT JOIN UserRoleAssignment ura ON u.id = ura.user_id
		LEFT JOIN RolePermission rp ON ura.role_id = rp.role_id
		LEFT JOIN Permission p ON rp.permission_id = p.id
		ORDER BY u.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usersMap := make(map[string]*UserWithPermissions)
	var userOrder []string
	for rows.Next() {
		var uid, email, firstName, lastName, phone, profilePic, lang sql.NullString
		var pid, pcode, pdesc sql.NullString
		err := rows.Scan(&uid, &email, &firstName, &lastName, &phone, &profilePic, &lang, &pid, &pcode, &pdesc)
		if err != nil {
			return nil, err
		}
		id := uid.String
		user, exists := usersMap[id]
		if !exists {
			user = &UserWithPermissions{
				ID:                id,
				Email:             email.String,
				FirstName:         firstName.String,
				LastName:          lastName.String,
				Phone:             phone.String,
				ProfilePictureURL: profilePic.String,
				PreferredLanguage: lang.String,
				Permissions:       []Permission{},
			}
			usersMap[id] = user
			userOrder = append(userOrder, id)
		}
		if pid.Valid {
			user.Permissions = append(user.Permissions, Permission{
				ID:          pid.String,
				Code:        pcode.String,
				Description: pdesc.String,
			})
		}
	}

	var result []UserWithPermissions
	for _, id := range userOrder {
		result = append(result, *usersMap[id])
	}
	return result, nil
}

// ChangeUserRole removes all roles from a user and assigns a new one by role_code
func (s *AuthService) ChangeUserRole(userID, roleCode string) error {
	var roleID string
	err := s.DB.QueryRow(`SELECT id FROM Role WHERE role_code = $1`, roleCode).Scan(&roleID)
	if err != nil {
		return err
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()
	// Remove all roles
	_, err = tx.Exec(`DELETE FROM UserRoleAssignment WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	// Assign new role by id
	_, err = tx.Exec(`INSERT INTO UserRoleAssignment (id, user_id, role_id, assigned_at) VALUES (gen_random_uuid(), $1, $2, NOW())`, userID, roleID)
	return err
}

// RemoveUserRole removes a specific role from a user by role_code
func (s *AuthService) RemoveUserRole(userID, roleCode string) error {
	var roleID string
	err := s.DB.QueryRow(`SELECT id FROM Role WHERE role_code = $1`, roleCode).Scan(&roleID)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`DELETE FROM UserRoleAssignment WHERE user_id = $1 AND role_id = $2`, userID, roleID)
	return err
}
