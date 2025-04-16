// internal/models/request/object.go
package request

// PostObjectRequest represents a POST request for retrieving an object
// This struct supports object requests with Passport authorization
// Used for:
// 1. Receiving request parameters in the POST /objects/{object_id} endpoint
// 2. Supporting Passport-based authorization and expanded object information
//
// Field descriptions:
// - Expand: Whether to return expanded object information, such as access method details (optional)
// - Passports: List of Passport tokens for authorization validation (optional)
type PostObjectRequest struct {
	Expand    bool     `json:"expand"`
	Passports []string `json:"passports"`
}

// PostAccessURLRequest represents a POST request for retrieving an access URL
// This struct supports access URL requests with Passport authorization
// Used for:
// 1. Receiving request parameters in the POST /objects/{object_id}/access/{access_id} endpoint
// 2. Supporting Passport-based authorization validation
//
// Field descriptions:
// - Passports: List of Passport tokens for authorization validation (optional)
type PostAccessURLRequest struct {
	Passports []string `json:"passports"`
}

// PostBulkAccessURLRequest represents a request for bulk access URL retrieval
// This struct extends the GA4GH DRS API specification to support bulk operations and Passport authorization
// Used for:
// 1. Receiving request parameters in the POST /objects/access and POST /objects/access/passports endpoints
// 2. Supporting retrieval of multiple access URLs for multiple objects in a single request
// 3. Supporting Passport-based authorization validation
//
// Field descriptions:
// - Passports: List of Passport tokens for authorization validation (optional)
// - BulkObjectAccessIDs: List of bulk object access IDs (required)
//   * ObjectID: Object ID (required)
//   * AccessIDs: List of access method IDs (required)
type PostBulkAccessURLRequest struct {
	Passports []string `json:"passports,omitempty"`
	BulkObjectAccessIDs []struct {
		ObjectID  string   `json:"object_id"`
		AccessIDs []string `json:"access_ids"`
	} `json:"bulk_object_access_ids"`
}
