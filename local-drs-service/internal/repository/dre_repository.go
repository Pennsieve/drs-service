// internal/repository/drs_repository.go
package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/daniellong/local-drs-service/internal/models"
)

// DrsRepository defines the interface for DRS data access
// This interface provides methods for retrieving object information and access URLs
type DrsRepository interface {
	GetObjectBasicInfo(ctx context.Context, id string) (*models.DrsObject, error)
	GetObjectAuthInfo(ctx context.Context, id string) (*models.DrsObject, error)
	//GetObjectData(ctx context.Context, id string) (io.ReadCloser, error)
	GetAccessURL(ctx context.Context, objectID, accessID string) (*models.AccessURL, error)
}

// LocalDrsRepository implements a local file system storage for DRS objects
// This implementation stores objects and their metadata in the local file system

type LocalDrsRepository struct {
	storageDir string                      // Directory for storing object data
	metaDir    string                      // Directory for storing object metadata
	objects    map[string]*models.DrsObject // In-memory cache of objects
	mutex      sync.RWMutex                // Mutex for thread-safe access to objects
}

// NewLocalDrsRepository creates a new instance of the local storage repository
// It initializes the storage and metadata directories and returns a repository instance
func NewLocalDrsRepository(storageDir string) (*LocalDrsRepository, error) {
	metaDir := filepath.Join(storageDir, "metadata")

	// Ensure directories exist
	for _, dir := range []string{storageDir, metaDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return &LocalDrsRepository{
		storageDir: storageDir,
		metaDir:    metaDir,
		objects:    make(map[string]*models.DrsObject),
	}, nil
}

// GetObjectBasicInfo retrieves basic object information for the GetObject service
// It returns the object's metadata without authorization-related fields
func (r *LocalDrsRepository) GetObjectBasicInfo(ctx context.Context, id string) (*models.DrsObject, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	obj, exists := r.objects[id]
	if !exists {
		return nil, models.NewNotFoundError(fmt.Sprintf("object not found: %s", id))
	}

	// Return basic information, excluding authorization-related fields
	basicObj := &models.DrsObject{
		ID:            obj.ID,
		DrsObjectID:   obj.DrsObjectID, // Ensure DrsObjectID matches ID
		Name:          obj.Name,
		SelfURI:       obj.SelfURI,
		Size:          obj.Size,
		CreatedTime:   obj.CreatedTime,
		UpdatedTime:   obj.UpdatedTime,
		Version:       obj.Version,
		MimeType:      obj.MimeType,
		Checksums:     obj.Checksums,
		AccessMethods: obj.AccessMethods,
		Content:       obj.Content,
		Description:   obj.Description,
		Aliases:       obj.Aliases,
	}

	return basicObj, nil
}

// GetObjectAuthInfo retrieves authorization information for the GetObjectAuthorization service
// It returns only the authorization-related fields of the object
func (r *LocalDrsRepository) GetObjectAuthInfo(ctx context.Context, id string) (*models.DrsObject, error) {
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

// GetObjectData retrieves object data, but is currently not used.
/*
func (r *LocalDrsRepository) GetObjectData(ctx context.Context, id string) (io.ReadCloser, error) {
	obj, err := r.GetObjectBasicInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(r.storageDir, obj.ID)
	return os.Open(filePath)
}
*/

// GetAccessURL retrieves an access URL for the specified object and access ID
// It checks if the object supports the "None" access type and returns the appropriate URL
func (r *LocalDrsRepository) GetAccessURL(ctx context.Context, objectID, accessID string) (*models.AccessURL, error) {
	obj, err := r.GetObjectBasicInfo(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Check if the object's SupportedTypes includes "None"
	hasNoneType := false
	for _, supportedType := range obj.SupportedTypes {
		if supportedType == "None" {
			hasNoneType = true
			break
		}
	}

	if !hasNoneType {
		// If the object doesn't support "None" type, return an empty URL indicating retry is needed
		return &models.AccessURL{
			URL: "",
		}, nil
	}

	for _, method := range obj.AccessMethods {
		if method.AccessID == accessID {
			// More complex logic could be implemented here, such as generating signed URLs
			accessURL := &models.AccessURL{
				URL:     method.AccessURL.URL,
				Headers: method.AccessURL.Headers,
			}

			return accessURL, nil
		}
	}

	return nil, fmt.Errorf("access ID not found")
}
