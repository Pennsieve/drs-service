// internal/models/object.go
package models

import "time"

// DrsObject represents a data object in the GA4GH DRS API specification
// This struct implements the DrsObject model from GA4GH DRS API v1.4.0
// Used for:
// 1. Returning object metadata in the GET /objects/{object_id} endpoint
// 2. As part of responses in bulk operations
// 3. As an internal data model for passing between service and repository layers
//
// Field descriptions:
// - ID: Internal unique identifier for storage and retrieval (not in spec, internal use)
// - DrsObjectID: Spec-compliant object ID used in API responses (spec requirement: required)
// - Name: Human-readable name for the object (spec requirement: optional)
// - SelfURI: Full URI of the object including service base URL (spec requirement: required)
// - Size: Size of the object in bytes (spec requirement: required)
// - CreatedTime: Time when the object was created (spec requirement: required)
// - UpdatedTime: Time when the object was last updated (spec requirement: optional)
// - Version: Object version identifier (spec requirement: optional)
// - MimeType: MIME type of the object (spec requirement: optional)
// - Checksums: List of checksums for the object, at least one required (spec requirement: required)
// - AccessMethods: List of methods to access the object (spec requirement: optional, but typically required in practice)
// - Content: Content description for composite objects (spec requirement: optional)
// - Description: Detailed description of the object (spec requirement: optional)
// - Aliases: List of alternative names for the object (spec requirement: optional)
// - SupportedTypes: List of supported access types, used for authorization checks (not in spec, internal use)
// - PassportAuthIssuers: List of supported Passport authorization issuers (not in spec, internal use)
// - BearerAuthIssuers: List of supported Bearer token issuers (not in spec, internal use)
type DrsObject struct {
	ID                  int            `json:"id"`
	DrsObjectID         string         `json:"drs_object_id"`
	Name                string         `json:"name,omitempty"`
	SelfURI             string         `json:"self_uri"`
	Size                int64          `json:"size"`
	CreatedTime         time.Time      `json:"created_time"`
	UpdatedTime         time.Time      `json:"updated_time,omitempty"`
	Version             string         `json:"version,omitempty"`
	MimeType            string         `json:"mime_type,omitempty"`
	Checksums           []Checksum     `json:"checksums"`
	AccessMethods       []AccessMethod `json:"access_methods,omitempty"`
	Content             []Content      `json:"content,omitempty"`
	Description         string         `json:"description,omitempty"`
	Aliases             []string       `json:"aliases,omitempty"`
	SupportedTypes      []string       `json:"supported_types"`
	PassportAuthIssuers []string       `json:"passport_auth_issuers,omitempty"`
	BearerAuthIssuers   []string       `json:"bearer_auth_issuers,omitempty"`
}

// Checksum represents checksum information for data integrity verification
// Used for:
// 1. As part of DrsObject to provide data integrity verification
// 2. Required field in the GA4GH DRS API specification
//
// Field descriptions:
// - Checksum: The checksum value (spec requirement: required)
// - Type: The checksum type, such as md5, sha256, etc. (spec requirement: required)
type Checksum struct {
	Checksum string `json:"checksum"`
	Type     string `json:"type"`
}

// Content represents content information for an object
// Used for:
// 1. Describing the content structure of composite objects
// 2. Optional field in the GA4GH DRS API specification
//
// Field descriptions:
// - Type: Content type (spec requirement: required)
// - Data: Content data (spec requirement: required)
type Content struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
