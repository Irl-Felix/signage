

package auth


import (
	   "log"
	   "errors"
	   "fmt"
	   "os"
	   "strings"
	   "time"

	   "github.com/MicahParks/keyfunc"
	   "github.com/golang-jwt/jwt/v4"
	   "github.com/Irl-Felix/signage/shared/util"
)

var jwks *keyfunc.JWKS

func init() {
	   projectURL := os.Getenv("SUPABASE_PROJECT_URL")
	   if projectURL == "" {
			   log.Println("SUPABASE_PROJECT_URL not set in environment")
			   panic("SUPABASE_PROJECT_URL not set in environment")
	   }

	   // Ensure no trailing slash
	   projectURL = strings.TrimSuffix(projectURL, "/")

	   jwksURL := fmt.Sprintf("%s/auth/v1/.well-known/jwks.json", projectURL)

	   var err error
	   jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
			   RefreshInterval: time.Hour,
			   RefreshErrorHandler: func(err error) {
					   util.LogError(err, "JWKS refresh error")
			   },
	   })
	   if err != nil {
			   util.LogError(err, "Failed to fetch JWKS from Supabase")
			   panic(fmt.Sprintf("Failed to fetch JWKS from Supabase: %v", err))
	   }
}

func VerifyJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, jwks.Keyfunc)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}
	if !token.Valid {
		return "", errors.New("invalid token: not valid")
	}

	// Check the algorithm in the header
	if token.Header["alg"] == nil {
		return "", errors.New("missing alg in token header")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok || strings.TrimSpace(sub) == "" {
		return "", errors.New("missing sub claim")
	}

	return sub, nil
}
