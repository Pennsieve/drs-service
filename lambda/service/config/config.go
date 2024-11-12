package config

import (
	"os"
)

type Config struct {
    DRSServiceID      string
    DRSOrgURL         string
    DocumentationURL  string
    CreatedAt         string
    UpdatedAt         string
    Environment       string
}

func NewConfig() Config {
    return Config{
        DRSServiceID:      getEnvOrDefault("DRS_SERVICE_ID", "net.pennsieve.drs"),
        DRSOrgURL:         getEnvOrDefault("DRS_ORG_URL", "https://pennsieve.dev"),
        DocumentationURL:  getEnvOrDefault("DOCUMENTATION_URL", "https://docs.pennsieve.io"),
        CreatedAt:         getEnvOrDefault("CREATED_AT", "2024-09-30T00:00:00Z"),
        UpdatedAt:         getEnvOrDefault("UPDATED_AT", "2024-09-30T00:00:00Z"),
        Environment:       getEnvOrDefault("ENVIRONMENT", "dev"),
    }
}

func getEnvOrDefault(key string, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
