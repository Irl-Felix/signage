package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserID      contextKey = "user_id"
	ContextPermissions contextKey = "permissions"
)

func Middleware(requiredPermission string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				log.Printf("[AUTH] Missing or malformed Authorization header on %s %s", r.Method, r.URL.Path)
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			user, permissions, err := ValidateTokenAndPermissions(token)
			if err != nil {
				log.Printf("[AUTH] Token validation error on %s %s: %v", r.Method, r.URL.Path, err)
					   if err == ErrNoRolesForUser {
							   log.Printf("[AUTH] No roles/permissions assigned to user %s on %s %s", user.Email, r.Method, r.URL.Path)
							   InsertAuditLog(user.ID, "PERMISSION_DENIED", fmt.Sprintf("No roles/permissions for %s on %s", requiredPermission, r.URL.Path))
							   http.Error(w, "Role not selected for user", http.StatusForbidden)
							   return
					   }
					   http.Error(w, "Invalid token", http.StatusUnauthorized)
					   return
			}

			if user == nil {
				log.Printf("[AUTH] User not found for valid token on %s %s", r.Method, r.URL.Path)
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			if len(permissions) == 0 {
				log.Printf("[AUTH] No roles/permissions assigned to user %s on %s %s", user.Email, r.Method, r.URL.Path)
				InsertAuditLog(user.ID, "PERMISSION_DENIED", fmt.Sprintf("No roles/permissions for %s on %s", requiredPermission, r.URL.Path))
				http.Error(w, "Role not selected for user", http.StatusForbidden)
				return
			}

			if requiredPermission != "" && !permissions[requiredPermission] {
				log.Printf("[AUTH] User %s lacks permission %s on %s %s", user.Email, requiredPermission, r.Method, r.URL.Path)
				InsertAuditLog(user.ID, "PERMISSION_DENIED", fmt.Sprintf("Denied %s on %s", requiredPermission, r.URL.Path))
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			InsertAuditLog(user.ID, "PERMISSION_GRANTED", fmt.Sprintf("Granted %s on %s", requiredPermission, r.URL.Path))
			log.Printf("[AUTH] User %s granted %s on %s %s", user.Email, requiredPermission, r.Method, r.URL.Path)

			ctx := context.WithValue(r.Context(), ContextUserID, user.ID)
			ctx = context.WithValue(ctx, ContextPermissions, permissions)
			next(w, r.WithContext(ctx))
		}
	}
}
