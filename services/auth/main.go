package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Irl-Felix/signage/services/auth/internal"
	"github.com/Irl-Felix/signage/shared/auth"
	"github.com/Irl-Felix/signage/shared/util"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {

	// Initialize shared auth DB connection for middleware
	if err := auth.InitDB(); err != nil {
		util.LogError(err, "Failed to init shared auth DB")
		panic("Failed to init shared auth DB")
	}

	// Connect to PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL not set, using default for local development")

	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		util.LogError(err, "Failed to connect to DB")
		panic("Failed to connect to DB")
	}
	defer db.Close()

	service := &internal.AuthService{DB: db}

	handler := &internal.Handler{Service: service}

	// ==============================
	// Public Authentication Endpoints
	// ==============================
	http.HandleFunc("POST /register", handler.Register)
	http.HandleFunc("POST /login", handler.Login)
	http.HandleFunc("POST /logout", handler.Logout)
	http.HandleFunc("POST /password-reset", handler.PasswordResetHandler)
	http.HandleFunc("POST /resend-verification", handler.ResendVerificationHandler)

	// ==============================
	// Admin Endpoints (Require \"MANAGE_USERS\" Permission)
	// ==============================

	   // --- User Management ---
	   http.HandleFunc("GET /admin/users", (handler.ListUsers))    // List all users
	   http.HandleFunc("GET /admin/users/{id}", auth.Middleware("MANAGE_USERS")(handler.GetUser)) // Get a specific user by ID
	   http.HandleFunc("POST /admin/users/assign-role", (handler.AssignUserRole))                 // Assign a role to a user
	   http.HandleFunc("GET /admin/users/{id}/permissions", (handler.GetUserPermissions))         // Get permissions for a specific user
	   http.HandleFunc("DELETE /admin/users/{id}", auth.Middleware("MANAGE_USERS")(handler.DeleteUser))
	   http.HandleFunc("GET /admin/users/stats", (handler.UserStatsHandler)) // User stats endpoint

	// --- Get user info with role and status ---
	// http.HandleFunc("GET /admin/users/{id}/info", auth.Middleware("MANAGE_USERS")(handler.GetUserInfo)) // Get user info with role and status

	http.HandleFunc("GET /admin/users/roles", (handler.ListUserRoles))                  // List all roles assigned to users
	http.HandleFunc("GET /admin/users/permissions", (handler.ListUsersWithPermissions)) // List all users with their permissions

	http.HandleFunc("POST /admin/users/changeuserrole", (handler.ChangeUserRole))                                    // Change user's role(As System Admin)
	http.HandleFunc("POST /admin/users/deleterolefromuser", auth.Middleware("MANAGE_USERS")(handler.RemoveUserRole)) // Remove a role from a user

	// --- Role Management ---
	http.HandleFunc("GET /admin/roles", auth.Middleware("MANAGE_USERS")(handler.ListRoles))
	http.HandleFunc("POST /admin/createroles", auth.Middleware("MANAGE_USERS")(handler.CreateRole))
	http.HandleFunc("GET /admin/role", auth.Middleware("MANAGE_USERS")(handler.GetRole))
	http.HandleFunc("POST /admin/updateroles", auth.Middleware("MANAGE_USERS")(handler.UpdateRole))
	http.HandleFunc("DELETE /admin/delroles", auth.Middleware("MANAGE_USERS")(handler.DeleteRole))

	// --- Permission Management ---
	http.HandleFunc("GET /admin/permissions", auth.Middleware("MANAGE_USERS")(handler.ListPermissions))
	http.HandleFunc("POST /admin/createpermissions", auth.Middleware("MANAGE_USERS")(handler.CreatePermission))
	http.HandleFunc("GET /admin/permission", auth.Middleware("MANAGE_USERS")(handler.GetPermission))
	http.HandleFunc("POST /admin/updatepermission", auth.Middleware("MANAGE_USERS")(handler.UpdatePermission))
	http.HandleFunc("DELETE /admin/delpermissions", auth.Middleware("MANAGE_USERS")(handler.DeletePermission))

	// --- Role-Permission Assignment ---
	http.HandleFunc("GET /admin/roles/getpermissions", auth.Middleware("MANAGE_USERS")(handler.GetRolePermissions))
	http.HandleFunc("POST /admin/roles/assignpermissions", auth.Middleware("MANAGE_USERS")(handler.AssignPermissionsToRole))
	http.HandleFunc("DELETE /admin/roles/removepermissions", auth.Middleware("MANAGE_USERS")(handler.RemovePermissionFromRole))

	// --- Audit Logs ---
	http.HandleFunc("GET /admin/audit-logs", auth.Middleware("MANAGE_USERS")(handler.ListAuditLogs))

	// ==============================
	// Dummy endpoints for testing
	// ==============================
	http.HandleFunc("GET /protected", auth.Middleware("MANAGE_USERS")(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "You have access to the protected endpoint!",
		})
	}))

	http.HandleFunc("GET /public", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "This is a public endpoint!",
		})
	})

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		util.LogError(err, "Server failed")
		panic("Server failed")
	}
}
