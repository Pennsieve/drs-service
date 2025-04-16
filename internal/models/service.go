// internal/models/service.go
package models

import "time"

// ServiceInfo represents metadata information for the DRS service
// This struct implements the Service Info model from GA4GH DRS API v1.4.0
// Used for:
// 1. Returning service metadata in the GET /service-info endpoint
// 2. Providing information about service capabilities and configuration to clients
//
// Field descriptions:
// - ID: Unique identifier for the service (spec requirement: required)
// - Name: Name of the service (spec requirement: required)
// - Type: Service type information (spec requirement: required)
// - Description: Detailed description of the service (spec requirement: optional)
// - Organization: Information about the organization providing the service (spec requirement: required)
// - ContactURL: URL for contacting the service provider (spec requirement: optional)
// - DocumentationURL: URL for service documentation (spec requirement: optional)
// - CreatedAt: Time when the service was created (spec requirement: optional)
// - UpdatedAt: Time when the service was last updated (spec requirement: optional)
// - Environment: Service runtime environment, e.g., prod, test, dev (spec requirement: optional)
// - Version: Service version (spec requirement: required)
// - MaxBulkRequestLength: Maximum number of objects supported in bulk requests (not in spec, extension field)
type ServiceInfo struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	Type                 ServiceType  `json:"type"`
	Description          string       `json:"description,omitempty"`
	Organization         Organization `json:"organization"`
	ContactURL           string       `json:"contactUrl,omitempty"`
	DocumentationURL     string       `json:"documentationUrl,omitempty"`
	CreatedAt            time.Time    `json:"createdAt,omitempty"`
	UpdatedAt            time.Time    `json:"updatedAt,omitempty"`
	Environment          string       `json:"environment,omitempty"`
	Version              string       `json:"version"`
	MaxBulkRequestLength int          `json:"maxBulkRequestLength"`
}

// ServiceType represents service type information
// This struct implements the ServiceType model from GA4GH Service Info specification
// Used for:
// 1. As part of ServiceInfo to describe the type and version of the service
// 2. Helping clients identify service compatibility
//
// Field descriptions:
// - Group: Service group name, typically the reversed organization domain (spec requirement: required)
// - Artifact: Service name (spec requirement: required)
// - Version: Service API version (spec requirement: required)
type ServiceType struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
}

// Organization represents organization information
// This struct implements the Organization model from GA4GH Service Info specification
// Used for:
// 1. As part of ServiceInfo to describe the organization providing the service
// 2. Providing information about service ownership and responsibility
//
// Field descriptions:
// - Name: Organization name (spec requirement: required)
// - URL: Organization website URL (spec requirement: optional, not included in this implementation)
type Organization struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}
