package auth

import "fmt"

import "errors"

var ErrNoRolesForUser = errors.New("no roles found for user")

func ValidateTokenAndPermissions(token string) (*UserProfile, map[string]bool, error) {
	// Verify token with Supabase (JWT validation)
	userID, err := VerifyJWT(token)
	if err != nil {
		fmt.Printf("Token verification failed: %v\n", err)
		return nil, nil, err
	}

	// Get user from DB by Supabase UID
	user, err := GetUserBySupabaseUID(userID)
	if err != nil {
		fmt.Printf("Failed to get user by Supabase UID %s: %v\n", userID, err)
		return nil, nil, err
	}

	// Get all roles for the user
	   roles, err := GetRolesByUserID(user.ID)
	   if err != nil {
			   if err.Error() == "no roles found for user" {
					   return user, nil, ErrNoRolesForUser
			   }
			   fmt.Printf("Failed to get roles for user %s: %v\n", user.ID, err)
			   return nil, nil, err
	   }

	// Aggregate permissions from all roles
	permMap := make(map[string]bool)
	for _, role := range roles {
		perms, err := GetPermissionsByRoleID(role.ID)
		fmt.Printf("Permissions for role %s: %+v\n", role.ID, perms)
		if err != nil {
			continue // skip roles with no permissions
		}
		for _, p := range perms {
			permMap[p] = true
		}
	}

	if len(permMap) == 0 {
		return user, permMap, nil // user exists but has no permissions
	}

	return user, permMap, nil
}
