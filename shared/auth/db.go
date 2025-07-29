package auth

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func InitDB() error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	db = conn
	return nil
}

func GetUserBySupabaseUID(uid string) (*UserProfile, error) {
	query := `SELECT id, supabase_uid, email, first_name, last_name, profile_picture_url, preferred_language, created_at FROM UserProfile WHERE supabase_uid = $1`
	row := db.QueryRow(context.Background(), query, uid)

	var user UserProfile
	err := row.Scan(&user.ID, &user.SupabaseUID, &user.Email, &user.FirstName, &user.LastName, &user.ProfilePictureURL, &user.PreferredLanguage, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetRolesByUserID fetches all roles for a user via UserRoleAssignment and Role tables
func GetRolesByUserID(userID string) ([]Role, error) {
	query := `
		SELECT r.id, r.role_code, r.name, r.scope
		FROM UserRoleAssignment ura
		JOIN Role r ON ura.role_id = r.id
		WHERE ura.user_id = $1
	`
	rows, err := db.Query(context.Background(), query, userID)
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
	if len(roles) == 0 {
		return nil, errors.New("no roles found for user")
	}
	return roles, nil
}

func GetPermissionsByRoleID(roleID string) ([]string, error) {
	query := `
		SELECT p.code
		FROM RolePermission rp
		JOIN Permission p ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`

	rows, err := db.Query(context.Background(), query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}

	if len(perms) == 0 {
		return nil, errors.New("no permissions found")
	}
	return perms, nil
}

func InsertAuditLog(userID, action, details string) {
	_, _ = db.Exec(context.Background(),
		`INSERT INTO AuditLog (user_id, action, details) VALUES ($1, $2, $3)`,
		userID, action, details,
	)
}
