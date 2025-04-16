// internal/models/response/bulk.go
package response

import (
	model "github.com/daniellong/local-drs-service/internal/models"
)

// BulkObjectResponse represents a response for bulk object retrieval
// This struct extends the GA4GH DRS API specification to support bulk operations
// Used for:
// 1. In the response of POST /objects endpoint
// 2. Supporting partially successful bulk requests with detailed success/failure information
//
// Field descriptions:
// - Summary: Request processing summary (required)
//   * Requested: Total number of requested objects
//   * Resolved: Number of successfully resolved objects
//   * Unresolved: Number of unresolved objects
// - ResolvedObjects: List of successfully resolved DRS objects (required)
// - UnresolvedObjects: Information about unresolved objects, grouped by error code (optional)
//   * ErrorCode: Error code, such as 404 for not found, 403 for forbidden, etc.
//   * ObjectIDs: List of object IDs with the same error
// - RetryAfter: Suggested time (in seconds) for the client to retry, used for partially successful requests (optional)
type BulkObjectResponse struct {
	Summary struct {
		Requested  int `json:"requested"`
		Resolved   int `json:"resolved"`
		Unresolved int `json:"unresolved"`
	} `json:"summary"`
	
	ResolvedObjects   []*model.DrsObject `json:"resolved_drs_object"`
	UnresolvedObjects []struct {
		ErrorCode int      `json:"error_code"`
		ObjectIDs []string `json:"object_ids"`
	} `json:"unresolved_drs_objects,omitempty"`
	RetryAfter *int `json:"retry_after,omitempty"`
}

// BulkAccessURLResponse represents a response for bulk access URL retrieval
// This struct extends the GA4GH DRS API specification to support bulk operations
// Used for:
// 1. In the response of POST /objects/access endpoint
// 2. Supporting partially successful bulk requests with detailed success/failure information
//
// Field descriptions:
// - Summary: Request processing summary (required)
//   * Requested: Total number of requested access URLs
//   * Resolved: Number of successfully resolved access URLs
//   * Unresolved: Number of unresolved access URLs
// - AccessURLs: List of successfully resolved access URLs (required)
// - Unresolved: Information about unresolved access URLs, grouped by error code (optional)
//   * ErrorCode: Error code, such as 404 for not found, 403 for forbidden, etc.
//   * ObjectIDs: List of object IDs with the same error
type BulkAccessURLResponse struct {
	Summary struct {
		Requested  int `json:"requested"`
		Resolved   int `json:"resolved"`
		Unresolved int `json:"unresolved"`
	} `json:"summary"`
	
	AccessURLs []model.BulkAccessURL `json:"access_urls"`
	Unresolved []struct {
		ErrorCode int      `json:"error_code"`
		ObjectIDs []string `json:"object_ids"`
	} `json:"unresolved,omitempty"`
}
