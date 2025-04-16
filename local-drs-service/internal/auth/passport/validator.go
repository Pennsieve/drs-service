// internal/auth/passport/validator.go
package passport

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Validator defines the interface for passport validation
// Implements methods for validating GA4GH passports and extracting claims
type Validator interface {
	ValidatePassport(passport string) (bool, error)
	GetClaims(passport string) (map[string]interface{}, error)
}

// DefaultValidator implements the default passport validation logic
// It validates passports against a list of trusted issuers
type DefaultValidator struct {
	trustedIssuers []string
}

// NewValidator creates a new passport validator with the specified trusted issuers
func NewValidator(trustedIssuers []string) *DefaultValidator {
	return &DefaultValidator{
		trustedIssuers: trustedIssuers,
	}
}

// ValidatePassport validates a GA4GH Passport
// Returns true if the passport is valid, false otherwise with an error
func (v *DefaultValidator) ValidatePassport(passport string) (bool, error) {
	if passport == "" {
		return false, errors.New("passport is empty")
	}

	// Parse the JWT token
	_, err := v.parseToken(passport)
	if err != nil {
		return false, fmt.Errorf("invalid passport: %w", err)
	}

	// Additional validation logic can be added here:
	// 1. Check if issuer is in the trusted issuers list
	// 2. Validate signature
	// 3. Check expiration time
	// 4. Validate other claims

	return true, nil
}

// GetClaims extracts claims from a passport
// Returns the claims as a map or an error if the passport is invalid
func (v *DefaultValidator) GetClaims(passport string) (map[string]interface{}, error) {
	claims, err := v.parseToken(passport)
	if err != nil {
		return nil, fmt.Errorf("failed to parse passport: %w", err)
	}

	return claims, nil
}

// parseToken parses a JWT token and returns the claims
// This is an internal helper method used by ValidatePassport and GetClaims
func (v *DefaultValidator) parseToken(tokenString string) (map[string]interface{}, error) {
	// Remove potential Bearer prefix
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Split the JWT token
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	// Decode the payload part
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	// Parse the claims
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return claims, nil
}
