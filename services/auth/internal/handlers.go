package internal

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Irl-Felix/signage/services/auth/internal/client"
	"github.com/Irl-Felix/signage/shared/util"
)

type Handler struct {
	Service *AuthService
}

// ==============================
// Public Authentication Handlers
// ==============================

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		FirstName         string `json:"first_name"`
		LastName          string `json:"last_name"`
		Phone             string `json:"phone"`
		ProfilePictureURL string `json:"profile_picture_url"`
		PreferredLanguage string `json:"preferred_language"`
	}
	if err := util.DecodeJSONBody(w, r, &body); err != nil {
		return
	}

	session, err := client.SignUp(body.Email, body.Password)
	if err != nil {
		util.HandleError(w, err, "Failed to sign up user", http.StatusBadRequest)
		return
	}

	if session.ID == "" || session.Email == "" {
		util.HandleError(w, nil, "Invalid Supabase signup response", http.StatusInternalServerError)
		return
	}

	user := UserProfile{
		SupabaseUID:       session.ID,
		Email:             session.Email,
		FirstName:         body.FirstName,
		LastName:          body.LastName,
		Phone:             body.Phone,
		ProfilePictureURL: body.ProfilePictureURL,
		PreferredLanguage: body.PreferredLanguage,
	}

	err = h.Service.CreateUser(&user)
	if err != nil {
		util.HandleError(w, err, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered",
		"user_id": user.ID,
		"email":   user.Email,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := util.DecodeJSONBody(w, r, &body); err != nil {
		return
	}

	session, err := client.Login(body.Email, body.Password)
	if err != nil {
		util.HandleError(w, err, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if session.User.ID == "" || session.User.Email == "" {
		util.HandleError(w, nil, "Invalid response from Supabase", http.StatusInternalServerError)
		return
	}

	user, err := h.Service.GetUserBySupabaseUID(session.User.ID)
	if err != nil {
		util.HandleError(w, err, "User profile not found", http.StatusNotFound)
		return
	}

	var tokenExpiresAt *time.Time
	if session.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(session.ExpiresIn) * time.Second)
		tokenExpiresAt = &t
	}
	sessionLog := SessionLog{
		UserID:         user.ID,
		IPAddress:      r.RemoteAddr,
		UserAgent:      r.UserAgent(),
		AccessToken:    session.AccessToken,
		RefreshToken:   session.RefreshToken,
		TokenExpiresAt: tokenExpiresAt,
		IsActive:       true,
	}

	err = h.Service.CreateSession(&sessionLog)
	if err != nil {
		util.HandleError(w, err, "Failed to log session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Login successful",
		"user_id":       user.ID,
		"email":         user.Email,
		"access_token":  session.AccessToken,
		"refresh_token": session.RefreshToken,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		util.HandleError(w, nil, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	accessToken := authHeader[7:]

	var tokenExpiresAt *time.Time
	var isActive bool
	err := h.Service.DB.QueryRow("SELECT token_expires_at, is_active FROM SessionLog WHERE access_token = $1 ORDER BY login_at DESC LIMIT 1", accessToken).Scan(&tokenExpiresAt, &isActive)
	if err != nil {
		util.HandleError(w, err, "Session not found", http.StatusUnauthorized)
		return
	}
	if tokenExpiresAt != nil && tokenExpiresAt.Before(time.Now()) {
		if isActive {
			_ = h.Service.LogoutSession(accessToken, *tokenExpiresAt)
		}
		util.HandleError(w, nil, "Session already expired", http.StatusUnauthorized)
		return
	}

	supabaseErr := client.Logout(accessToken)
	if supabaseErr != nil {
		util.HandleError(w, supabaseErr, "Failed to logout from Supabase", http.StatusInternalServerError)
		return
	}

	dbErr := h.Service.LogoutSession(accessToken, time.Now())
	if dbErr != nil {
		util.HandleError(w, dbErr, "Logged out from Supabase, but failed to update session log", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

func (h *Handler) PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := util.DecodeJSONBody(w, r, &body); err != nil {
		return
	}
	err := client.RequestPasswordReset(body.Email)
	if err != nil {
		util.HandleError(w, err, "Failed to send password reset email", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset email sent if the address exists."})
}

func (h *Handler) ResendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := util.DecodeJSONBody(w, r, &body); err != nil {
		return
	}
	err := client.ResendVerificationEmail(body.Email)
	if err != nil {
		util.HandleError(w, err, "Failed to resend verification email", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Verification email resent if the address exists."})
}

// ==============================
// User Profile Management
// ==============================

func (h *Handler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID                string `json:"id"`
		FirstName         string `json:"first_name"`
		LastName          string `json:"last_name"`
		Email             string `json:"email"`
		Phone             string `json:"phone"`
		ProfilePictureURL string `json:"profile_picture_url"`
		PreferredLanguage string `json:"preferred_language"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	err := h.Service.UpdateUserProfile(&UserProfile{
		ID:                input.ID,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Email:             input.Email,
		Phone:             input.Phone,
		ProfilePictureURL: input.ProfilePictureURL,
		PreferredLanguage: input.PreferredLanguage,
	})
	if err != nil {
		util.HandleError(w, err, "Failed to update user profile", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User profile updated",
	})
}

// ==============================
// Admin API Handlers
// ==============================

// --- User Management ---

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	   users, err := h.Service.ListUsers()
	   if err != nil {
			   util.HandleError(w, err, "Failed to list users", http.StatusInternalServerError)
			   return
	   }
	   w.Header().Set("Content-Type", "application/json")
	   json.NewEncoder(w).Encode(map[string]interface{}{
			   "users": users,
	   })
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Expecting /admin/users/:id, extract id from URL
	   id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	   var (
			   user UserProfile
			   profilePictureURL, preferredLanguage sql.NullString
	   )
	   err := h.Service.DB.QueryRow(`SELECT id, supabase_uid, first_name, last_name, email, phone, profile_picture_url, preferred_language, created_at FROM UserProfile WHERE id = $1`, id).Scan(
			   &user.ID, &user.SupabaseUID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &profilePictureURL, &preferredLanguage, &user.CreatedAt)
	   if err != nil {
			   util.HandleError(w, err, "User not found", http.StatusNotFound)
			   return
	   }
	   user.ProfilePictureURL = profilePictureURL.String
	   user.PreferredLanguage = preferredLanguage.String
	   w.Header().Set("Content-Type", "application/json")
	   json.NewEncoder(w).Encode(user)
}

func (h *Handler) AssignUserRole(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		RoleCode string `json:"role_code"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if err := h.Service.AddUserRole(input.UserID, input.RoleCode); err != nil {
		util.HandleError(w, err, "Failed to assign role", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User assigned to role",
	})
}

func (h *Handler) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	// Expecting /admin/users/:id/permissions, extract user id from URL
	id := r.URL.Path[strings.Index(r.URL.Path, "/admin/users/")+13:]
	if idx := strings.Index(id, "/"); idx != -1 {
		id = id[:idx]
	}
	// Get all roles for the user
	rows, err := h.Service.DB.Query(`
			   SELECT p.id, p.code, p.description
			   FROM UserRoleAssignment ura
			   JOIN RolePermission rp ON ura.role_id = rp.role_id
			   JOIN Permission p ON rp.permission_id = p.id
			   WHERE ura.user_id = $1`, id)
	if err != nil {
		util.HandleError(w, err, "Failed to get user permissions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type Permission struct {
		ID          string `json:"id"`
		Code        string `json:"code"`
		Description string `json:"description"`
	}
	var perms []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Description); err != nil {
			util.HandleError(w, err, "Failed to scan permission", http.StatusInternalServerError)
			return
		}
		perms = append(perms, p)
	}
	w.Header().Set("Content-Type", "application/json")
	if len(perms) == 0 {
		// Log a clear message for server logs
		log.Printf("[AUTH] User %s has role(s) but no permissions mapped. Check RolePermission table for role assignments.", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No permissions assigned to this user. The user's role(s) exist but may not have any permissions mapped.",
		})
		return
	}
	json.NewEncoder(w).Encode(perms)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Expecting /admin/users/:id, extract id from URL
	id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	res, err := h.Service.DB.Exec(`DELETE FROM UserProfile WHERE id = $1`, id)
	if err != nil {
		util.HandleError(w, err, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		util.HandleError(w, nil, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}

// --- User-Role Assignment ---

func (h *Handler) ListUserRoles(w http.ResponseWriter, r *http.Request) {
	assignments, err := h.Service.ListUserRoles()
	if err != nil {
		util.HandleError(w, err, "Failed to list user-role assignments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assignments)
}

// --- User Permissions Overview ---

// ListUsersWithPermissions returns all users and their permissions
func (h *Handler) ListUsersWithPermissions(w http.ResponseWriter, r *http.Request) {
	// Query all users and their permissions
	rows, err := h.Service.DB.Query(`
		SELECT u.id, u.email, u.first_name, u.last_name, u.phone, u.profile_picture_url, u.preferred_language,
			   p.id, p.code, p.description
		FROM UserProfile u
		LEFT JOIN UserRoleAssignment ura ON u.id = ura.user_id
		LEFT JOIN RolePermission rp ON ura.role_id = rp.role_id
		LEFT JOIN Permission p ON rp.permission_id = p.id
		ORDER BY u.id
	`)
	if err != nil {
		util.HandleError(w, err, "Failed to list users with permissions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Permission struct {
		ID          string `json:"id"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
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

	usersMap := make(map[string]*UserWithPermissions)
	var userOrder []string
	for rows.Next() {
		var uid, email, firstName, lastName, phone, profilePic, lang sql.NullString
		var pid, pcode, pdesc sql.NullString
		err := rows.Scan(&uid, &email, &firstName, &lastName, &phone, &profilePic, &lang, &pid, &pcode, &pdesc)
		if err != nil {
			util.HandleError(w, err, "Failed to scan user-permission row", http.StatusInternalServerError)
			return
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

	// Build result in order
	var result []UserWithPermissions
	for _, id := range userOrder {
		result = append(result, *usersMap[id])
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// --- Audit Logs ---

func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.Service.ListAuditLogs()
	if err != nil {
		util.HandleError(w, err, "Failed to list audit logs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// --- Role Management ---

func (h *Handler) ListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.Service.ListRoles()
	if err != nil {
		util.HandleError(w, err, "Failed to list roles", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func (h *Handler) GetRole(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"role_id": "..."}
	var input struct {
		RoleID string `json:"role_id"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.RoleID == "" {
		util.HandleError(w, nil, "role_id is required", http.StatusBadRequest)
		return
	}
	var role Role
	err := h.Service.DB.QueryRow(`SELECT id, role_code, name, scope FROM Role WHERE id = $1`, input.RoleID).Scan(&role.ID, &role.RoleCode, &role.Name, &role.Scope)
	if err != nil {
		util.HandleError(w, err, "Role not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RoleCode string `json:"role_code"`
		Name     string `json:"name"`
		Scope    string `json:"scope"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	role := &Role{
		RoleCode: input.RoleCode,
		Name:     input.Name,
		Scope:    input.Scope,
	}
	if err := h.Service.CreateRole(role); err != nil {
		util.HandleError(w, err, "Failed to create role", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Role created",
	})
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"role_id": "...", "role_code": "...", "name": "...", "scope": "..."}
	var input struct {
		RoleID   string `json:"role_id"`
		RoleCode string `json:"role_code"`
		Name     string `json:"name"`
		Scope    string `json:"scope"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.RoleID == "" {
		util.HandleError(w, nil, "role_id is required", http.StatusBadRequest)
		return
	}
	res, err := h.Service.DB.Exec(`UPDATE Role SET role_code = $1, name = $2, scope = $3 WHERE id = $4`, input.RoleCode, input.Name, input.Scope, input.RoleID)
	if err != nil {
		util.HandleError(w, err, "Failed to update role", http.StatusInternalServerError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		util.HandleError(w, nil, "Role not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Role updated"})
}

func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"role_id": "..."}
	var input struct {
		RoleID string `json:"role_id"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.RoleID == "" {
		util.HandleError(w, nil, "role_id is required", http.StatusBadRequest)
		return
	}
	res, err := h.Service.DB.Exec(`DELETE FROM Role WHERE id = $1`, input.RoleID)
	if err != nil {
		util.HandleError(w, err, "Failed to delete role", http.StatusInternalServerError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		util.HandleError(w, nil, "Role not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Role deleted"})
}

// --- Permission Management ---

func (h *Handler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Service.DB.Query(`SELECT id, code, description FROM Permission`)
	if err != nil {
		util.HandleError(w, err, "Failed to list permissions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type Permission struct {
		ID          string `json:"id"`
		Code        string `json:"code"`
		Description string `json:"description"`
	}
	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Description); err != nil {
			util.HandleError(w, err, "Failed to scan permission", http.StatusInternalServerError)
			return
		}
		permissions = append(permissions, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

func (h *Handler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	var id string
	err := h.Service.DB.QueryRow(
		`INSERT INTO Permission (id, code, description) VALUES (gen_random_uuid(), $1, $2) RETURNING id`,
		input.Code, input.Description,
	).Scan(&id)
	if err != nil {
		util.HandleError(w, err, "Failed to create permission", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Permission created",
		"id":      id,
	})
}

func (h *Handler) GetPermission(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"permission_id": "..."}
	var input struct {
		PermissionID string `json:"permission_id"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.PermissionID == "" {
		util.HandleError(w, nil, "permission_id is required", http.StatusBadRequest)
		return
	}
	var perm struct {
		ID          string `json:"id"`
		Code        string `json:"code"`
		Description string `json:"description"`
	}
	err := h.Service.DB.QueryRow(`SELECT id, code, description FROM Permission WHERE id = $1`, input.PermissionID).Scan(&perm.ID, &perm.Code, &perm.Description)
	if err != nil {
		util.HandleError(w, err, "Permission not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perm)
}

func (h *Handler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"permission_id": "...", "code": "...", "name": "...", "description": "..."}
	var input struct {
		PermissionID string `json:"permission_id"`
		Code         string `json:"code"`
		Description  string `json:"description"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.PermissionID == "" {
		util.HandleError(w, nil, "permission_id is required", http.StatusBadRequest)
		return
	}
	res, err := h.Service.DB.Exec(
		`UPDATE Permission SET code = $1, description = $2 WHERE id = $3`,
		input.Code, input.Description, input.PermissionID,
	)
	if err != nil {
		util.HandleError(w, err, "Failed to update permission", http.StatusInternalServerError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		util.HandleError(w, nil, "Permission not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Permission updated"})
}

func (h *Handler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"permission_id": "..."}
	var input struct {
		PermissionID string `json:"permission_id"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.PermissionID == "" {
		util.HandleError(w, nil, "permission_id is required", http.StatusBadRequest)
		return
	}
	res, err := h.Service.DB.Exec(`DELETE FROM Permission WHERE id = $1`, input.PermissionID)
	if err != nil {
		util.HandleError(w, err, "Failed to delete permission", http.StatusInternalServerError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		util.HandleError(w, nil, "Permission not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Permission deleted"})
}

// --- Linkin Role-Permission Assignment ---

func (h *Handler) ListRolePermissions(w http.ResponseWriter, r *http.Request) {
	mappings, err := h.Service.ListRolePermissions()
	if err != nil {
		util.HandleError(w, err, "Failed to list role-permission mappings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mappings)
}

func (h *Handler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"role_id": "..."}
	var input struct {
		RoleID string `json:"role_id"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.RoleID == "" {
		util.HandleError(w, nil, "role_id is required", http.StatusBadRequest)
		return
	}
	perms, err := h.Service.GetRolePermissions(input.RoleID)
	if err != nil {
		util.HandleError(w, err, "Failed to get role permissions", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perms)
}

func (h *Handler) AssignPermissionsToRole(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"role_id": "...", "permission_ids": ["...", ...]}
	var input struct {
		RoleID        string   `json:"role_id"`
		PermissionIDs []string `json:"permission_ids"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.RoleID == "" {
		util.HandleError(w, nil, "role_id is required", http.StatusBadRequest)
		return
	}
	if err := h.Service.AssignPermissionsToRole(input.RoleID, input.PermissionIDs); err != nil {
		util.HandleError(w, err, "Failed to assign permissions to role", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Permissions assigned to role"})
}

func (h *Handler) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	// Accepts JSON body: {"role_id": "...", "permission_id": "..."}
	var input struct {
		RoleID       string `json:"role_id"`
		PermissionID string `json:"permission_id"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if input.RoleID == "" || input.PermissionID == "" {
		util.HandleError(w, nil, "role_id and permission_id are required", http.StatusBadRequest)
		return
	}
	ok, err := h.Service.RemovePermissionFromRole(input.RoleID, input.PermissionID)
	if err != nil {
		util.HandleError(w, err, "Failed to remove permission from role", http.StatusInternalServerError)
		return
	}
	if !ok {
		util.HandleError(w, nil, "Role or permission not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Permission removed from role"})
}

// ChangeUserRole changes a user's role (removes all roles and assigns a new one)
func (h *Handler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		RoleCode string `json:"role_code"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if err := h.Service.ChangeUserRole(input.UserID, input.RoleCode); err != nil {
		util.HandleError(w, err, "Failed to change user role", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User role changed",
	})
}

// RemoveUserRole removes a specific role from a user
func (h *Handler) RemoveUserRole(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string `json:"user_id"`
		RoleCode string `json:"role_code"`
	}
	if err := util.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if err := h.Service.RemoveUserRole(input.UserID, input.RoleCode); err != nil {
		util.HandleError(w, err, "Failed to remove user role", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User role removed",
	})
}

// UserStatsHandler returns user statistics (counts by status)
func (h *Handler) UserStatsHandler(w http.ResponseWriter, r *http.Request) {
	   stats, err := h.Service.GetUserStats()
	   if err != nil {
			   util.HandleError(w, err, "Failed to get user stats", http.StatusInternalServerError)
			   return
	   }
	   w.Header().Set("Content-Type", "application/json")
	   json.NewEncoder(w).Encode(stats)
}