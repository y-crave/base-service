package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type jwtClaimsMinimal struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
}

// claimsFromBearerJWTUnsafe decodes JWT payload without verifying the signature.
// For local development only when Istio RequestAuthentication is absent.
func claimsFromBearerJWTUnsafe(bearer string) (jwtClaimsMinimal, error) {
	token := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(bearer), "Bearer "))
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return jwtClaimsMinimal{}, errors.New("jwt: expected three segments")
	}
	payload := parts[1]
	raw, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		raw, err = base64URLEncodingDecodePadding(payload)
	}
	if err != nil {
		return jwtClaimsMinimal{}, err
	}
	var c jwtClaimsMinimal
	if err := json.Unmarshal(raw, &c); err != nil {
		return jwtClaimsMinimal{}, err
	}
	return c, nil
}

func base64URLEncodingDecodePadding(s string) ([]byte, error) {
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	return base64.URLEncoding.DecodeString(s)
}
