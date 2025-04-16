// internal/models/access.go
package models

// AccessMethod represents a method to access a data object
// This struct implements the AccessMethod model from GA4GH DRS API v1.4.0
// Used for:
// 1. As part of DrsObject to describe how to access the object data
// 2. In the GET /objects/{object_id}/access/{access_id} endpoint for generating access URLs
//
// Field descriptions:
// - Type: Access type such as file, https, s3, gs, etc. (spec requirement: required)
// - AccessURL: Direct access URL, can be null (spec requirement: either AccessURL or AccessID must be present)
// - AccessID: Unique identifier for the access method, used to retrieve access URL later (spec requirement: either AccessID or AccessURL must be present)
// - Region: Geographic region where the data is located, for selecting the nearest data replica (spec requirement: optional)
type AccessMethod struct {
	Type      string     `json:"type"`
	AccessURL *AccessURL `json:"access_url,omitempty"`
	AccessID  string     `json:"access_id,omitempty"`
	Region    string     `json:"region,omitempty"`
}

// AccessURL represents an access URL for an object
// This struct implements the AccessURL model from GA4GH DRS API v1.4.0
// Used for:
// 1. As part of AccessMethod to provide direct access URL
// 2. In the response of GET /objects/{object_id}/access/{access_id} endpoint
//
// Field descriptions:
// - URL: URL to access the object data (spec requirement: required)
// - Headers: HTTP headers needed when accessing the URL (spec requirement: optional)
//   Using map[string]string instead of fixed structure provides more flexible header support
type AccessURL struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

// BulkAccessURL represents a single access URL in a bulk access URL response
// This struct extends the GA4GH DRS API specification for bulk operations
// Used for:
// 1. In bulk access URL responses as a single object's access URL
// 2. Supports requests with or without Passport
//
// Field descriptions:
// - DrsObjectID: Unique identifier for the object (required)
// - DrsAccessID: Unique identifier for the access method (required)
// - URL: URL to access the object data (required)
// - Headers: HTTP headers needed when accessing the URL (optional)
//   Using map[string]string instead of fixed structure provides more flexible header support
type BulkAccessURL struct {
	DrsObjectID string            `json:"drs_object_id"`
	DrsAccessID string            `json:"drs_access_id"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// Access type constants
// These constants define the supported access types, compliant with GA4GH DRS API specification
const (
	AccessTypeFile  = "file"  // Local file access
	AccessTypeHTTPS = "https" // HTTPS URL access
	AccessTypeS3    = "s3"    // Amazon S3 access
	AccessTypeGS    = "gs"    // Google Cloud Storage access
)
