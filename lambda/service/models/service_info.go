package models

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

func NewServiceInfo(id, url, documentationURL, createdAt, updatedAt, environment string) ServiceInfo {
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
		DocumentationURL: documentationURL,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		Environment:      environment,
		Version:          "1.0.0",
	}
}
