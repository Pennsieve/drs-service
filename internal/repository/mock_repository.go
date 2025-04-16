// internal/repository/mock_repository.go
package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pennsieve/drs-service/internal/models"
)

// MockDrsRepository implements DrsRepository interface with in-memory test data
// This implementation is used for testing and demonstration purposes
type MockDrsRepository struct {
	objects map[string]*models.DrsObject // In-memory storage of objects
	mutex   sync.RWMutex                 // Mutex for thread-safe access to objects
}

// NewMockDrsRepository creates a new instance of the mock repository
// It initializes the repository with predefined test data
func NewMockDrsRepository() (*MockDrsRepository, error) {
	repo := &MockDrsRepository{
		objects: make(map[string]*models.DrsObject),
	}

	// Initialize with test data
	repo.initializeTestData()

	return repo, nil
}

// initializeTestData populates the repository with test data
func (r *MockDrsRepository) initializeTestData() {
	// Create test objects
	testObjects := []models.DrsObject{
		{
			ID:          1,
			DrsObjectID: "test-object-1",
			Name:        "Test Object 1",
			SelfURI:     "drs://test.example.com/test-object-1",
			Size:        1024,
			CreatedTime: time.Now().Add(-24 * time.Hour),
			UpdatedTime: time.Now(),
			Version:     "1.0",
			MimeType:    "application/json",
			Checksums: []models.Checksum{
				{
					Checksum: "d41d8cd98f00b204e9800998ecf8427e",
					Type:     "md5",
				},
			},
			AccessMethods: []models.AccessMethod{
				{
					Type:     models.AccessTypeHTTPS,
					AccessID: "access-1",
					AccessURL: &models.AccessURL{
						URL: "https://test.example.com/data/test-object-1",
						Headers: map[string]string{
							"Authorization": "Bearer test-token",
						},
					},
				},
			},
			Description:    "A test object for demonstration",
			Aliases:        []string{"test1", "demo1"},
			SupportedTypes: []string{"None"}, // No authorization required
		},
		{
			ID:          2,
			DrsObjectID: "test-object-2",
			Name:        "Test Object 2",
			SelfURI:     "drs://test.example.com/test-object-2",
			Size:        2048,
			CreatedTime: time.Now().Add(-48 * time.Hour),
			UpdatedTime: time.Now(),
			Version:     "1.0",
			MimeType:    "application/octet-stream",
			Checksums: []models.Checksum{
				{
					Checksum: "e99a18c428cb38d5f260853678922e03",
					Type:     "md5",
				},
			},
			AccessMethods: []models.AccessMethod{
				{
					Type:     models.AccessTypeHTTPS,
					AccessID: "access-2",
					AccessURL: &models.AccessURL{
						URL: "https://test.example.com/data/test-object-2",
					},
				},
			},
			Description:         "A protected test object requiring authorization",
			Aliases:             []string{"test2", "demo2"},
			SupportedTypes:      []string{"Passport"}, // Requires passport authorization
			PassportAuthIssuers: []string{"https://auth.example.com"},
		},
		{
			ID:          3,
			DrsObjectID: "test-object-3",
			Name:        "Test Object 3",
			SelfURI:     "drs://test.example.com/test-object-3",
			Size:        4096,
			CreatedTime: time.Now().Add(-72 * time.Hour),
			UpdatedTime: time.Now(),
			Version:     "1.0",
			MimeType:    "text/plain",
			Checksums: []models.Checksum{
				{
					Checksum: "900150983cd24fb0d6963f7d28e17f72",
					Type:     "md5",
				},
			},
			AccessMethods: []models.AccessMethod{
				{
					Type:     models.AccessTypeHTTPS,
					AccessID: "access-3",
				},
			},
			Description:       "A test object with delayed access URL",
			Aliases:           []string{"test3", "demo3"},
			SupportedTypes:    []string{"None"}, // No authorization required
			BearerAuthIssuers: []string{"https://bearer.example.com"},
		},
	}

	// Add test objects to the repository
	for _, obj := range testObjects {
		objCopy := obj // Create a copy to avoid reference issues
		r.objects[objCopy.DrsObjectID] = &objCopy
	}
}

// GetObjectBasicInfo retrieves basic object information for the GetObject service
func (r *MockDrsRepository) GetObjectBasicInfo(ctx context.Context, id string) (*models.DrsObject, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	obj, exists := r.objects[id]
	if !exists {
		return nil, models.NewNotFoundError(fmt.Sprintf("object not found: %s", id))
	}

	// Return basic information, excluding authorization-related fields
	basicObj := &models.DrsObject{
		ID:             obj.ID,
		DrsObjectID:    obj.DrsObjectID,
		Name:           obj.Name,
		SelfURI:        obj.SelfURI,
		Size:           obj.Size,
		CreatedTime:    obj.CreatedTime,
		UpdatedTime:    obj.UpdatedTime,
		Version:        obj.Version,
		MimeType:       obj.MimeType,
		Checksums:      obj.Checksums,
		AccessMethods:  obj.AccessMethods,
		Content:        obj.Content,
		Description:    obj.Description,
		Aliases:        obj.Aliases,
		SupportedTypes: obj.SupportedTypes, // Include SupportedTypes for authorization checks
	}

	return basicObj, nil
}

// GetObjectAuthInfo retrieves authorization information for the GetObjectAuthorization service
func (r *MockDrsRepository) GetObjectAuthInfo(ctx context.Context, id string) (*models.DrsObject, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	obj, exists := r.objects[id]
	if !exists {
		return nil, models.NewNotFoundError(fmt.Sprintf("object not found: %s", id))
	}

	// Return authorization-related information
	authObj := &models.DrsObject{
		DrsObjectID:         obj.DrsObjectID,
		SupportedTypes:      obj.SupportedTypes,
		PassportAuthIssuers: obj.PassportAuthIssuers,
		BearerAuthIssuers:   obj.BearerAuthIssuers,
	}

	return authObj, nil
}

// GetAccessURL retrieves an access URL for the specified object and access ID
func (r *MockDrsRepository) GetAccessURL(ctx context.Context, objectID, accessID string) (*models.AccessURL, error) {
	obj, err := r.GetObjectBasicInfo(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Special case for test-object-3: simulate a delayed access URL
	if objectID == "test-object-3" {
		return &models.AccessURL{
			URL: "", // Empty URL indicates the resource is not yet ready
		}, nil
	}

	for _, method := range obj.AccessMethods {
		if method.AccessID == accessID {
			if method.AccessURL != nil {
				return method.AccessURL, nil
			}

			// If AccessURL is not directly available, create a mock one
			return &models.AccessURL{
				URL:     fmt.Sprintf("https://test.example.com/data/%s", objectID),
				Headers: map[string]string{"X-Test": "true"},
			}, nil
		}
	}

	return nil, models.NewNotFoundError(fmt.Sprintf("access ID not found: %s", accessID))
}
