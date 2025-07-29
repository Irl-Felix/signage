package internal

import "time"

type UserProfile struct {
	ID                string // UUID
	SupabaseUID       string
	FirstName         string
	LastName          string
	Email             string
	Phone             string
	ProfilePictureURL string
	PreferredLanguage string
	CreatedAt         time.Time
}

type SessionLog struct {
	ID             string // UUID
	UserID         string
	LoginAt        time.Time
	LogoutAt       *time.Time
	IPAddress      string
	UserAgent      string
	AccessToken    string
	RefreshToken   string
	TokenExpiresAt *time.Time
	IsActive       bool
	DurationSec    *int
}
type Role struct {
	ID       string
	RoleCode string
	Name     string
	Scope    string
}

type RolePermission struct {
	ID           string
	RoleID       string
	PermissionID string
}
