package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Irl-Felix/signage/shared/util"
)

var (
	SUPABASE_API_URL = os.Getenv("SUPABASE_PROJECT_URL")
	API_KEY          = os.Getenv("SUPABASE_ANON_KEY")
)

// For /register (sign-up)
type SupabaseSignupResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// For /login
type SupabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type SupabaseLoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         SupabaseUser `json:"user"`
	ExpiresIn    int          `json:"expires_in"` // seconds until expiry
}

func SignUp(email, password string) (*SupabaseSignupResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	var session SupabaseSignupResponse
	err := callSupabase("/auth/v1/signup", payload, &session)
	return &session, err
}

func Login(email, password string) (*SupabaseLoginResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	var session SupabaseLoginResponse
	err := callSupabase("/auth/v1/token?grant_type=password", payload, &session)
	return &session, err
}
func Logout(accessToken string) error {
	req, err := http.NewRequest("POST", SUPABASE_API_URL+"/auth/v1/logout", nil)
	if err != nil {
		return err
	}
	req.Header.Set("apikey", API_KEY)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase logout error: %s", string(body))
	}

	return nil
}

func callSupabase(endpoint string, payload map[string]string, result interface{}) error {
	jsonBody, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", SUPABASE_API_URL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		util.LogError(err, "Failed to create Supabase request")
		return err
	}
	req.Header.Set("apikey", API_KEY)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.LogError(err, "Failed to execute Supabase request")
		return err
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)
	// Optionally log response for debugging only on error

	if res.StatusCode >= 400 {
		util.LogError(fmt.Errorf("supabase error: %s", string(bodyBytes)), "Supabase returned error response")
		return fmt.Errorf("supabase error: %s", string(bodyBytes))
	}

	err = json.Unmarshal(bodyBytes, result)
	if err != nil {
		util.LogError(err, "Failed to decode Supabase response")
		return fmt.Errorf("decode error: %v", err)
	}

	return nil
}

func RequestPasswordReset(email string) error {
	payload := map[string]string{
		"email": email,
	}
	// Supabase expects no response body for this endpoint
	return callSupabaseNoResult("/auth/v1/recover", payload)
}

func ResendVerificationEmail(email string) error {
	payload := map[string]string{
		"email": email,
		"type":  "signup",
	}
	return callSupabaseNoResult("/auth/v1/resend", payload)
}

// dont chk for response body-unlike callSupabase.
func callSupabaseNoResult(endpoint string, payload map[string]string) error {
	jsonBody, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", SUPABASE_API_URL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		util.LogError(err, "Failed to create Supabase request")
		return err
	}
	req.Header.Set("apikey", API_KEY)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.LogError(err, "Failed to execute Supabase request")
		return err
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		util.LogError(fmt.Errorf("supabase error: %s", string(bodyBytes)), "Supabase returned error response")
		return fmt.Errorf("supabase error: %s", string(bodyBytes))
	}
	return nil
}
