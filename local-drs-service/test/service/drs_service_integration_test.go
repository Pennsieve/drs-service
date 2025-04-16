// test/service/drs_service_integration_test.go
package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daniellong/local-drs-service/config"
	"github.com/daniellong/local-drs-service/internal/handler"
	"github.com/daniellong/local-drs-service/internal/models"
	"github.com/daniellong/local-drs-service/internal/repository"
	"github.com/daniellong/local-drs-service/internal/router"
	"github.com/daniellong/local-drs-service/internal/service"
)

// setupTestServer creates a test HTTP server with the DRS API
func setupTestServer(t *testing.T) *httptest.Server {
	// Create mock repository
	repo, err := repository.NewMockDrsRepository()
	require.NoError(t, err)

	// Create test config
	cfg := config.LoadConfig()

	// Create service
	drsService := service.NewDrsService(repo, "https://test.example.com", cfg)

	// Create handler
	drsHandler := handler.NewDrsHandler(drsService)

	// Setup routes
	r := router.SetupRoutes(drsHandler)

	// Create test server
	return httptest.NewServer(r)
}

// TestServiceInfoEndpoint tests the /service-info endpoint
func TestServiceInfoEndpoint(t *testing.T) {
	// Setup
	server := setupTestServer(t)
	defer server.Close()

	// Execute
	resp, err := http.Get(fmt.Sprintf("%s/ga4gh/drs/v1/service-info", server.URL))

	// Verify
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var info models.ServiceInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify response contents
	assert.NotEmpty(t, info.ID)
	assert.NotEmpty(t, info.Name)
	assert.Equal(t, "1.4.0", info.Type.Version)
}

// TestGetObjectEndpoint tests the /objects/{object_id} endpoint
func TestGetObjectEndpoint(t *testing.T) {
	// Setup
	server := setupTestServer(t)
	defer server.Close()

	// Test cases
	testCases := []struct {
		name       string
		objectID   string
		statusCode int
	}{
		{
			name:       "Valid object",
			objectID:   "test-object-1",
			statusCode: http.StatusOK,
		},
		{
			name:       "Non-existent object",
			objectID:   "non-existent",
			statusCode: http.StatusNotFound,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make request
			resp, err := http.Get(fmt.Sprintf("%s/ga4gh/drs/v1/objects/%s", server.URL, tc.objectID))
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify status code
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			// If success, verify response contents
			if tc.statusCode == http.StatusOK {
				var obj models.DrsObject
				err = json.NewDecoder(resp.Body).Decode(&obj)
				require.NoError(t, err)

				assert.Equal(t, tc.objectID, obj.DrsObjectID)
				assert.NotEmpty(t, obj.Name)
			}
		})
	}
}

// TestGetAccessURLEndpoint tests the /objects/{object_id}/access/{access_id} endpoint
func TestGetAccessURLEndpoint(t *testing.T) {
	// Setup
	server := setupTestServer(t)
	defer server.Close()

	// Test cases
	testCases := []struct {
		name       string
		objectID   string
		accessID   string
		statusCode int
		retryAfter bool
	}{
		{
			name:       "Valid access URL",
			objectID:   "test-object-1",
			accessID:   "access-1",
			statusCode: http.StatusOK,
			retryAfter: false,
		},
		{
			name:       "Delayed access URL",
			objectID:   "test-object-3",
			accessID:   "access-3",
			statusCode: http.StatusAccepted,
			retryAfter: true,
		},
		{
			name:       "Non-existent object",
			objectID:   "non-existent",
			accessID:   "access-1",
			statusCode: http.StatusNotFound,
			retryAfter: false,
		},
	}

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make request
			resp, err := http.Get(fmt.Sprintf("%s/ga4gh/drs/v1/objects/%s/access/%s",
				server.URL, tc.objectID, tc.accessID))
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify status code
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			// Check Retry-After header if expected
			if tc.retryAfter {
				assert.NotEmpty(t, resp.Header.Get("Retry-After"))
			}

			// If success or accepted, verify response contents
			if tc.statusCode == http.StatusOK || tc.statusCode == http.StatusAccepted {
				var accessURL models.AccessURL
				err = json.NewDecoder(resp.Body).Decode(&accessURL)
				require.NoError(t, err)

				if tc.statusCode == http.StatusAccepted {
					assert.Empty(t, accessURL.URL)
				} else {
					assert.NotEmpty(t, accessURL.URL)
				}
			}
		})
	}
}

// TestBulkObjectsEndpoint tests the POST /objects endpoint
func TestBulkObjectsEndpoint(t *testing.T) {
	// Setup
	server := setupTestServer(t)
	defer server.Close()

	// Create request body
	reqBody := struct {
		Passports     []string `json:"passports"`
		BulkObjectIDs []string `json:"bulk_object_ids"`
	}{
		BulkObjectIDs: []string{"test-object-1", "test-object-2", "non-existent"},
	}

	// Serialize request
	jsonData, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/ga4gh/drs/v1/objects", server.URL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var bulkResp struct {
		Summary struct {
			Requested  int `json:"requested"`
			Resolved   int `json:"resolved"`
			Unresolved int `json:"unresolved"`
		} `json:"summary"`
		UnresolvedDrsObjects []struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		} `json:"unresolved_drs_objects"`
		ResolvedDrsObject []models.DrsObject `json:"resolved_drs_object"`
	}

	err = json.NewDecoder(resp.Body).Decode(&bulkResp)
	require.NoError(t, err)

	// Verify response contents
	assert.Equal(t, 3, bulkResp.Summary.Requested)

	// Note: test-object-2 requires Passport authorization, so without providing a valid Passport,
	// only test-object-1 will be resolved, while test-object-2 and non-existent will be in the unresolved list
	assert.Equal(t, 1, bulkResp.Summary.Resolved)
	assert.Equal(t, 2, bulkResp.Summary.Unresolved)

	// Verify resolved objects - only test-object-1 should be resolved
	assert.Len(t, bulkResp.ResolvedDrsObject, 1)

	// Verify unresolved objects - server will combine objects with the same error code
	assert.Len(t, bulkResp.UnresolvedDrsObjects, 1)
	assert.Equal(t, http.StatusNotFound, bulkResp.UnresolvedDrsObjects[0].ErrorCode)

	// Verify that the unresolved object ID list contains the expected IDs
	assert.Len(t, bulkResp.UnresolvedDrsObjects[0].ObjectIDs, 2)
	assert.Contains(t, bulkResp.UnresolvedDrsObjects[0].ObjectIDs, "test-object-2")
	assert.Contains(t, bulkResp.UnresolvedDrsObjects[0].ObjectIDs, "non-existent")

	// Verify object IDs in resolved objects
	var foundIDs []string
	for _, obj := range bulkResp.ResolvedDrsObject {
		foundIDs = append(foundIDs, obj.DrsObjectID)
	}
	assert.Contains(t, foundIDs, "test-object-1")
	// test-object-2 requires authorization, so it should not be in the resolved list
	assert.NotContains(t, foundIDs, "test-object-2")
}
