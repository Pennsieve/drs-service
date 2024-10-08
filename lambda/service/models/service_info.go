package models

import (
	"os"
)


type ServiceInfo struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Type             TypeInfo `json:"type"`
	Description      string   `json:"description"`
	Organization     OrgInfo  `json:"organization"`
	ContactURL       string   `json:"contactUrl"`
	DocumentationURL string   `json:"documentationUrl"`
	CreatedAt        string   `json:"createdAt"`
	UpdatedAt        string   `json:"updatedAt"`
	Environment      string   `json:"environment"`   
	Version          string   `json:"version"`       
}


type TypeInfo struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
}


type OrgInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}


func NewServiceInfo() ServiceInfo {

	id := os.Getenv("DRS_SERVICE_ID")
	if id == "" {
		id = "io.pennsieve.drs" // default
	}

	url := os.Getenv("DRS_ORG_URL")
	if url == "" {
		url = "https://pennsieve.io" // default
	}

	return ServiceInfo{
		ID:          id,
		Name:        "Pennsieve DRS Service",
		Type:        TypeInfo{Group: "org.ga4gh", Artifact: "drs"},
		Description: "This service provides an API that conforms to the GA4GH DRS specifications.",
		Organization: OrgInfo{
			Name: "Pennsieve",
			URL:  url,
		},
		ContactURL:       "support@pennsieve.io",
		DocumentationURL: "https://docs.pennsieve.io",
		CreatedAt:        "2024-09-30T00:00:00Z",
		UpdatedAt:        "2024-09-30T00:00:00Z",
		Environment:      "test",
		Version:          "1.0.0",
	}
}
