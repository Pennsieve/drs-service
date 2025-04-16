// internal/models/request/bulk.go
package request

// BulkObjectRequest represents a request for bulk object retrieval
// This struct extends the GA4GH DRS API specification to support bulk operations
// Used for:
// 1. Receiving request parameters in the POST /objects endpoint
// 2. Supporting retrieval of multiple objects in a single request
// 3. Supporting Passport-based authorization validation
//
// Field descriptions:
// - BulkObjectIDs: List of object IDs to retrieve (required)
// - Passports: List of Passport tokens for authorization validation (optional)
// - Expand: Whether to return expanded object information (optional)
type BulkObjectRequest struct {
	BulkObjectIDs []string `json:"bulk_object_ids"`
	Passports     []string `json:"passports,omitempty"`
	Expand        bool     `json:"expand"`
}

// BulkAccessRequest represents a request for bulk access URL retrieval
// This struct extends the GA4GH DRS API specification to support bulk operations
// Used for:
// 1. Receiving request parameters in the POST /objects/access endpoint
// 2. Supporting retrieval of multiple access URLs for multiple objects in a single request
// 3. Supporting Passport-based authorization validation
//
// Field descriptions:
// - Objects: List of objects with their access IDs (required)
//   * ObjectID: Object ID (required)
//   * AccessIDs: List of access method IDs for the object (required)
// - Passports: List of Passport tokens for authorization validation (optional)
type BulkAccessRequest struct {
	Objects []struct {
		ObjectID  string   `json:"object_id"`
		AccessIDs []string `json:"access_ids"`
	} `json:"objects"`
	Passports []string `json:"passports,omitempty"`
}
