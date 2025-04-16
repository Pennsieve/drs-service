// config/loader.go
package config

import (
	"os"
	"strconv"
	"strings"
)

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	config := NewDefaultConfig()
	
	// 服务配置
	if port, err := strconv.Atoi(getEnvOrDefault("SERVER_PORT", "")); err == nil && port > 0 {
		config.ServerPort = port
	}
	if baseURL := getEnvOrDefault("BASE_URL", ""); baseURL != "" {
		config.BaseURL = baseURL
	}
	
	// 存储配置
	if storageDir := getEnvOrDefault("STORAGE_DIR", ""); storageDir != "" {
		config.StorageDir = storageDir
	}
	
	// 服务信息
	if serviceID := getEnvOrDefault("DRS_SERVICE_ID", ""); serviceID != "" {
		config.DRSServiceID = serviceID
	}
	if serviceName := getEnvOrDefault("SERVICE_NAME", ""); serviceName != "" {
		config.ServiceName = serviceName
	}
	if serviceDesc := getEnvOrDefault("SERVICE_DESCRIPTION", ""); serviceDesc != "" {
		config.ServiceDescription = serviceDesc
	}
	if orgName := getEnvOrDefault("ORGANIZATION_NAME", ""); orgName != "" {
		config.OrganizationName = orgName
	}
	if orgURL := getEnvOrDefault("DRS_ORG_URL", ""); orgURL != "" {
		config.DRSOrgURL = orgURL
	}
	if docURL := getEnvOrDefault("DOCUMENTATION_URL", ""); docURL != "" {
		config.DocumentationURL = docURL
	}
	if createdAt := getEnvOrDefault("CREATED_AT", ""); createdAt != "" {
		config.CreatedAt = createdAt
	}
	if updatedAt := getEnvOrDefault("UPDATED_AT", ""); updatedAt != "" {
		config.UpdatedAt = updatedAt
	}
	if env := getEnvOrDefault("ENVIRONMENT", ""); env != "" {
		config.Environment = env
	}
	
	// 认证配置
	config.EnableBasicAuth = getEnvAsBool("ENABLE_BASIC_AUTH", config.EnableBasicAuth)
	config.EnableBearerAuth = getEnvAsBool("ENABLE_BEARER_AUTH", config.EnableBearerAuth)
	config.EnablePassportAuth = getEnvAsBool("ENABLE_PASSPORT_AUTH", config.EnablePassportAuth)
	
	if issuers := getEnvOrDefault("TRUSTED_ISSUERS", ""); issuers != "" {
		config.TrustedIssuers = strings.Split(issuers, ",")
	}
	
	// 批量操作配置
	if maxBulk, err := strconv.Atoi(getEnvOrDefault("MAX_BULK_REQUEST_LENGTH", "")); err == nil && maxBulk > 0 {
		config.MaxBulkRequestLength = maxBulk
	}
	
	return config
}

// 辅助函数
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}
