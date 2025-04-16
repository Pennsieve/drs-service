// test/service/drs_service_test.go
package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pennsieve/drs-service/internal/config"
	"github.com/pennsieve/drs-service/internal/repository"
	"github.com/pennsieve/drs-service/internal/service"
)

// setupTestService creates a new DRS service with mock repository for testing
func setupTestService(t *testing.T) service.DrsService {
	// Create mock repository
	repo, err := repository.NewMockDrsRepository()
	require.NoError(t, err)

	// Create test config
	cfg := &config.Config{
		ServerPort:           8080,
		BaseURL:              "https://test.example.com",
		DRSServiceID:         "test-drs-service",
		ServiceName:          "Test DRS Service",
		ServiceDescription:   "A test DRS service for unit testing",
		OrganizationName:     "Test Organization",
		DRSOrgURL:            "https://test.example.com",
		DocumentationURL:     "https://test.example.com/docs",
		CreatedAt:            time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		UpdatedAt:            time.Now().Format(time.RFC3339),
		Environment:          "test",
		MaxBulkRequestLength: 100,
	}

	// Create service
	return service.NewDrsService(repo, cfg.BaseURL, cfg)
}

// TestGetServiceInfo tests the GetServiceInfo method
func TestGetServiceInfo(t *testing.T) {
	// Setup
	svc := setupTestService(t)
	ctx := context.Background()

	// Execute
	info, err := svc.GetServiceInfo(ctx)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "test-drs-service", info.ID)
	assert.Equal(t, "Test DRS Service", info.Name)
	assert.Equal(t, "1.4.0", info.Type.Version)
}

// TestGetObject tests the GetObject method
func TestGetObject(t *testing.T) {
	// Setup
	svc := setupTestService(t)
	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name      string
		objectID  string
		passports []string
		expectErr bool
	}{
		{
			name:      "Valid object without auth",
			objectID:  "test-object-1",
			passports: nil,
			expectErr: false,
		},
		{
			name:      "Object requiring auth without passport",
			objectID:  "test-object-2",
			passports: nil,
			expectErr: true, // Should fail without passport
		},
		{
			name:      "Non-existent object",
			objectID:  "non-existent",
			passports: nil,
			expectErr: true,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			obj, err := svc.GetObject(ctx, tc.objectID, tc.passports)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, obj)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, obj)
				assert.Equal(t, tc.objectID, obj.DrsObjectID)
			}
		})
	}
}

// TestGetAccessURL tests the GetAccessURL method
func TestGetAccessURL(t *testing.T) {
	// Setup
	svc := setupTestService(t)
	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name        string
		objectID    string
		accessID    string
		expectErr   bool
		expectRetry bool
	}{
		{
			name:        "Valid access URL",
			objectID:    "test-object-1",
			accessID:    "access-1",
			expectErr:   false,
			expectRetry: false,
		},
		{
			name:        "Delayed access URL",
			objectID:    "test-object-3",
			accessID:    "access-3",
			expectErr:   false,
			expectRetry: true,
		},
		{
			name:        "Invalid object ID",
			objectID:    "non-existent",
			accessID:    "access-1",
			expectErr:   true,
			expectRetry: false,
		},
		{
			name:        "Invalid access ID",
			objectID:    "test-object-1",
			accessID:    "non-existent",
			expectErr:   true,
			expectRetry: false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, retryAfter, err := svc.GetAccessURL(ctx, tc.objectID, tc.accessID)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, url)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, url)

				if tc.expectRetry {
					assert.NotNil(t, retryAfter)
					assert.Empty(t, url.URL)
				} else {
					assert.Nil(t, retryAfter)
					assert.NotEmpty(t, url.URL)
				}
			}
		})
	}
}

// TestGetBulkObjects tests the GetBulkObjects method
func TestGetBulkObjects(t *testing.T) {
	// Setup
	svc := setupTestService(t)
	ctx := context.Background()

	// Execute
	response, retryAfter, err := svc.GetBulkObjects(ctx, []string{"test-object-1", "test-object-2", "non-existent"})

	// Verify
	assert.NoError(t, err)
	assert.Nil(t, retryAfter)
	assert.NotNil(t, response)
	
	// Verify summary
	assert.Equal(t, 3, response.Summary.Requested)
	assert.Equal(t, 2, response.Summary.Resolved)
	assert.Equal(t, 1, response.Summary.Unresolved)
	
	// Verify resolved objects
	assert.Len(t, response.ResolvedObjects, 2)
	
	// Verify object contents
	var foundIDs []string
	for _, obj := range response.ResolvedObjects {
		foundIDs = append(foundIDs, obj.DrsObjectID)
	}
	assert.Contains(t, foundIDs, "test-object-1")
	assert.Contains(t, foundIDs, "test-object-2")
	
	// Verify unresolved objects
	assert.Len(t, response.UnresolvedObjects, 1)
	assert.Equal(t, 404, response.UnresolvedObjects[0].ErrorCode)
	assert.Contains(t, response.UnresolvedObjects[0].ObjectIDs, "non-existent")
}
