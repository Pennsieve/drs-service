// config/config.go
package config

// Config 应用配置
type Config struct {
	// 服务配置
	ServerPort int
	BaseURL    string

	// 存储配置
	StorageDir string

	// 服务信息
	DRSServiceID       string
	ServiceName        string
	ServiceDescription string
	OrganizationName   string
	DRSOrgURL          string
	DocumentationURL   string
	CreatedAt          string
	UpdatedAt          string
	Environment        string

	// 认证配置
	EnableBasicAuth    bool
	EnableBearerAuth   bool
	EnablePassportAuth bool
	TrustedIssuers     []string

	// 批量操作配置
	MaxBulkRequestLength int
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	return &Config{
		ServerPort: 8080,
		BaseURL:    "localhost:8080",
		StorageDir: "./storage",

		DRSServiceID:       "net.pennsieve.drs",
		ServiceName:        "DRS Demo Service",
		ServiceDescription: "A demonstration implementation of GA4GH DRS",
		OrganizationName:   "DRS Demo Organization",
		DRSOrgURL:          "https://pennsieve.dev",
		DocumentationURL:   "https://docs.pennsieve.io",
		CreatedAt:          "2024-09-30T00:00:00Z",
		UpdatedAt:          "2024-09-30T00:00:00Z",
		Environment:        "dev",

		EnableBasicAuth:    true,
		EnableBearerAuth:   true,
		EnablePassportAuth: true,
		TrustedIssuers:     []string{},

		MaxBulkRequestLength: 100,
	}
}
