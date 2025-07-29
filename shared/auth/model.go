
package auth

import "time"


import "database/sql"

type UserProfile struct {
    ID                string         // UUID
    SupabaseUID       string
    FirstName         sql.NullString
    LastName          sql.NullString
    Email             string
    Phone             sql.NullString
    ProfilePictureURL sql.NullString
    PreferredLanguage sql.NullString
    CreatedAt         time.Time
}

// ToPlain returns a struct with all fields as string (empty if null)
func (u UserProfile) ToPlain() UserProfilePlain {
    return UserProfilePlain{
        ID:                u.ID,
        SupabaseUID:       u.SupabaseUID,
        FirstName:         nullToStr(u.FirstName),
        LastName:          nullToStr(u.LastName),
        Email:             u.Email,
        Phone:             nullToStr(u.Phone),
        ProfilePictureURL: nullToStr(u.ProfilePictureURL),
        PreferredLanguage: nullToStr(u.PreferredLanguage),
        CreatedAt:         u.CreatedAt,
    }
}

type UserProfilePlain struct {
    ID                string
    SupabaseUID       string
    FirstName         string
    LastName          string
    Email             string
    Phone             string
    ProfilePictureURL string
    PreferredLanguage string
    CreatedAt         time.Time
}

func nullToStr(ns sql.NullString) string {
    if ns.Valid {
        return ns.String
    }
    return ""
}
 

type Role struct {
    ID       string
    RoleCode string
    Name     string
    Scope    string // 'global' or 'business'
}

// UserRoleAssignment matches the join table for user-role assignments
type UserRoleAssignment struct {
    ID         string
    UserID     string
    RoleID     string
    BusinessID *string // nullable
    LocationID *string // nullable
    AssignedAt time.Time
}
type Permission struct {
    ID   string
    Name string
}

type RolePermission struct {
    RoleID       string
    PermissionID string
}
