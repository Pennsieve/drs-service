// internal/models/response/access.go
package response

import (
	model "github.com/pennsieve/drs-service/internal/models"
)

// AccessURLResponse represents an access URL response
// This struct encapsulates the AccessURL response from GA4GH DRS API v1.4.0
// Used for:
// 1. In the response of GET /objects/{object_id}/access/{access_id} endpoint
// 2. In the response of POST /objects/{object_id}/access/{access_id} endpoint
//
// Field descriptions:
// - AccessURL: Access information containing URL and optional Headers (spec requirement: required)
//   Implemented by embedding model.AccessURL to maintain consistency with the model layer
type AccessURLResponse struct {
	*model.AccessURL
}

// BulkAccessURLWithPassportResponse represents a bulk access URL response with Passport authorization
// This struct extends the GA4GH DRS API specification to support bulk operations and Passport authorization
// Used for:
// 1. In the response of POST /objects/access/passports endpoint
// 2. Supporting partially successful bulk requests with detailed success/failure information
//
// Field descriptions:
// - Summary: Request processing summary, including counts of requested, resolved, and unresolved objects (required)
//   * Requested: Total number of requested objects
//   * Resolved: Number of successfully resolved objects
//   * Unresolved: Number of unresolved objects
// - UnresolvedDrsObjects: Information about unresolved objects, grouped by error code (optional)
//   * ErrorCode: Error code, such as 404 for not found, 403 for forbidden, etc.
//   * ObjectIDs: List of object IDs with the same error
// - ResolvedDrsObjectAccessURLs: List of successfully resolved object access URLs (required)
// - RetryAfter: Suggested time (in seconds) for the client to retry, used for partially successful requests (optional)
//   When this field is set, the response status code should be 202 Accepted
type BulkAccessURLWithPassportResponse struct {
	Summary struct {
		Requested  int `json:"requested"`
		Resolved   int `json:"resolved"`
		Unresolved int `json:"unresolved"`
	} `json:"summary"`
	
	UnresolvedDrsObjects []struct {
		ErrorCode int      `json:"error_code"`
		ObjectIDs []string `json:"object_ids"`
	} `json:"unresolved_drs_objects,omitempty"`
	
	ResolvedDrsObjectAccessURLs []model.BulkAccessURL `json:"resolved_drs_object_access_urls"`
	
	RetryAfter *int `json:"retry_after,omitempty"`
}
