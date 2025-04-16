// internal/models/request/access.go
package request

// AccessURLRequest represents a request for access URL retrieval
// This struct supports access URL requests with optional Passport authorization
// Used for:
// 1. Receiving request parameters for access URL endpoints
// 2. Supporting Passport-based authorization validation
//
// Field description:
// - Passports: List of Passport tokens for authorization validation (optional)
type AccessURLRequest struct {
	Passports []string `json:"passports,omitempty"`
}
