// internal/service/drs_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/daniellong/local-drs-service/config"
	"github.com/daniellong/local-drs-service/internal/auth/passport"
	"github.com/daniellong/local-drs-service/internal/models"
	"github.com/daniellong/local-drs-service/internal/models/response"
	"github.com/daniellong/local-drs-service/internal/repository"
)

// DrsService defines the interface for DRS API operations
// This interface provides methods for all DRS API endpoints defined in the GA4GH specification
type DrsService interface {
	GetServiceInfo(ctx context.Context) (*models.ServiceInfo, error)
	GetObjectAuthorization(ctx context.Context, id string, passports []string) (*models.DrsObject, error)
	GetObject(ctx context.Context, id string, passports []string) (*models.DrsObject, error)
	GetObjectWithPassport(ctx context.Context, id string, expand bool, passports []string) (*models.DrsObject, error)
	GetBulkObjects(ctx context.Context, ids []string) (*response.BulkObjectResponse, *int, error)
	GetBulkObjectsWithPassport(ctx context.Context, ids []string, expand bool, passports []string) (*response.BulkObjectResponse, *int, error)
	GetAccessURL(ctx context.Context, objectID, accessID string) (*models.AccessURL, *int, error)
	GetAccessURLWithPassport(ctx context.Context, objectID, accessID string, passports []string) (*models.AccessURL, *int, error)
	GetBulkAccessURLs(ctx context.Context, objectAccessIDs map[string][]string) ([]models.BulkAccessURL, map[string]error)
	GetBulkAccessURLsWithPassport(ctx context.Context, objectAccessIDs map[string][]string, passports []string) (*response.BulkAccessURLWithPassportResponse, error)
	ListObjects(ctx context.Context, pageSize int, pageToken string) ([]models.DrsObject, string, error)
	GetBulkObjectAuthorizations(ctx context.Context, ids []string, passports []string) (*response.BulkObjectResponse, *int, error)
}

// DefaultDrsService implements the DRS service interface
// This implementation provides the business logic for all DRS API operations
type DefaultDrsService struct {
	repository repository.DrsRepository // Repository for data access
	baseURL    string                   // Base URL for the service
	config     *config.Config           // Service configuration
}

// NewDrsService creates a new DRS service instance
// It initializes the service with the provided repository, base URL, and configuration
func NewDrsService(repo repository.DrsRepository, baseURL string, cfg *config.Config) *DefaultDrsService {
	return &DefaultDrsService{
		repository: repo,
		baseURL:    baseURL,
		config:     cfg,
	}
}

// GetServiceInfo retrieves service metadata information
// This method implements the GET /service-info endpoint from the GA4GH DRS API specification

func (s *DefaultDrsService) GetServiceInfo(ctx context.Context) (*models.ServiceInfo, error) {
	// Parse time strings from configuration
	createdAt, err := time.Parse(time.RFC3339, s.config.CreatedAt)
	if err != nil {
		// If parsing fails, use a fixed default time
		defaultTime, _ := time.Parse(time.RFC3339, "2019-06-04T12:58:19Z")
		createdAt = defaultTime
	}

	updatedAt, err := time.Parse(time.RFC3339, s.config.UpdatedAt)
	if err != nil {
		// If parsing fails, use a fixed default time
		defaultTime, _ := time.Parse(time.RFC3339, "2019-06-04T12:58:19Z")
		updatedAt = defaultTime
	}

	// Build service information, including API basic information
	serviceInfo := &models.ServiceInfo{
		ID:   s.config.DRSServiceID,
		Name: s.config.ServiceName,
		Type: models.ServiceType{
			Group:    "org.ga4gh",
			Artifact: "drs",
			Version:  "1.4.0",
		},
		Description: s.config.ServiceDescription,
		Organization: models.Organization{
			Name: s.config.OrganizationName,
			URL:  s.config.DRSOrgURL,
		},
		DocumentationURL:     s.config.DocumentationURL,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
		Environment:          s.config.Environment,
		Version:              "0.1.0",
		MaxBulkRequestLength: s.config.MaxBulkRequestLength,
	}

	return serviceInfo, nil
}

// GetObject retrieves an object by its ID
// This method implements the GET /objects/{object_id} endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetObject(ctx context.Context, id string, passports []string) (*models.DrsObject, error) {
	// First, check if the object exists
	obj, err := s.repository.GetObjectBasicInfo(ctx, id)
	if err != nil {
		return nil, models.NewNotFoundError(fmt.Sprintf("Object %s does not exist", id))
	}

	// Check if the object requires authorization
	// If the object's SupportedTypes does not contain "None", then passport authorization is required
	requiresAuth := true
	if obj.SupportedTypes != nil {
		for _, t := range obj.SupportedTypes {
			if t == "None" {
				requiresAuth = false
				break
			}
		}
	} else {
		// If SupportedTypes is not specified, default to not requiring authorization
		requiresAuth = false
	}

	if requiresAuth {
		// If no passport is provided, deny access
		if len(passports) == 0 {
			return nil, models.NewUnauthorizedError("Passport authorization is required")
		}

		// Validate the passport
		authorized := false
		validator := passport.NewValidator([]string{})
		for _, p := range passports {
			valid, _ := validator.ValidatePassport(p)
			if valid {
				authorized = true
				break
			}
		}

		if !authorized {
			return nil, models.NewUnauthorizedError("Invalid passport")
		}
	}

	return obj, nil
}

// GetAccessURL retrieves an access URL for an object
// This method implements the GET /objects/{object_id}/access/{access_id} endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetAccessURL(ctx context.Context, objectID, accessID string) (*models.AccessURL, *int, error) {
	// First, get object basic information
	obj, err := s.repository.GetObjectBasicInfo(ctx, objectID)
	if err != nil {
		return nil, nil, models.NewNotFoundError("Object does not exist")
	}

	// Check if the object requires passport authorization
	requiresAuth := false
	for _, supportedType := range obj.SupportedTypes {
		if supportedType == "Passport" {
			requiresAuth = true
			break
		}
	}

	if requiresAuth {
		// Get auth info to provide more detailed error message
		authInfo, _ := s.repository.GetObjectAuthInfo(ctx, objectID)
		errorMsg := "Object requires Passport authorization"
		if len(authInfo.PassportAuthIssuers) > 0 {
			errorMsg += " from issuers: " + fmt.Sprint(authInfo.PassportAuthIssuers)
		}
		return nil, nil, models.NewUnauthorizedError(errorMsg)
	}

	accessURL, err := s.repository.GetAccessURL(ctx, objectID, accessID)
	if err != nil {
		return nil, nil, models.NewNotFoundError(fmt.Sprintf("Access ID %s does not exist", accessID))
	}

	// Check if the URL is ready
	// If the URL is empty, it means the resource is being prepared and needs to be retried later
	if accessURL.URL == "" {
		retryAfter := 60 // default 60 seconds to retry
		return accessURL, &retryAfter, nil
	}

	return accessURL, nil, nil
}

// GetAccessURLWithPassport retrieves an access URL for an object with passport authorization
// This method implements the GET /objects/{object_id}/access/{access_id} endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetAccessURLWithPassport(ctx context.Context, objectID, accessID string, passports []string) (*models.AccessURL, *int, error) {
	// First, get the object information
	obj, err := s.repository.GetObjectBasicInfo(ctx, objectID)
	if err != nil {
		return nil, nil, models.NewNotFoundError(fmt.Sprintf("Object %s does not exist", objectID))
	}

	// Check if the object requires passport authorization
	requiresAuth := false
	for _, supportedType := range obj.SupportedTypes {
		if supportedType == "Passport" {
			requiresAuth = true
			break
		}
	}

	if requiresAuth {
		// If no passport is provided, deny access
		if len(passports) == 0 {
			return nil, nil, models.NewUnauthorizedError("Passport authorization is required")
		}

		// Validate the passport
		authorized := false
		validator := passport.NewValidator(obj.PassportAuthIssuers)
		for _, p := range passports {
			valid, _ := validator.ValidatePassport(p)
			if valid {
				authorized = true
				break
			}
		}

		if !authorized {
			return nil, nil, models.NewUnauthorizedError("Invalid passport")
		}
	}

	// Get the access URL
	accessURL, err := s.repository.GetAccessURL(ctx, objectID, accessID)
	if err != nil {
		return nil, nil, models.NewNotFoundError(fmt.Sprintf("Access ID %s does not exist", accessID))
	}

	// Check if the URL is ready
	// If the URL is empty, it means the resource is being prepared and needs to be retried later
	if accessURL.URL == "" {
		retryAfter := 60 // default 60 seconds to retry
		return accessURL, &retryAfter, nil
	}

	return accessURL, nil, nil
}

// GetBulkObjects retrieves multiple objects by their IDs
// This method implements the POST /objects endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetBulkObjects(ctx context.Context, ids []string) (*response.BulkObjectResponse, *int, error) {
	// Create response structure
	bulkResponse := &response.BulkObjectResponse{
		Summary: struct {
			Requested  int `json:"requested"`
			Resolved   int `json:"resolved"`
			Unresolved int `json:"unresolved"`
		}{
			Requested:  len(ids),
			Resolved:   0,
			Unresolved: 0,
		},
		ResolvedObjects: make([]*models.DrsObject, 0),
		UnresolvedObjects: make([]struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		}, 0),
	}

	// Check if the request exceeds the maximum limit
	if len(ids) > s.config.MaxBulkRequestLength {
		return nil, nil, models.NewBadRequestError(fmt.Sprintf("Bulk request exceeds maximum limit %d", s.config.MaxBulkRequestLength))
	}

	// Group unresolved objects by error code
	unresolvedByErrorCode := make(map[int][]string)

	// Process each object ID
	for _, id := range ids {
		// Get the baseinfo for every object
		obj, err := s.repository.GetObjectBasicInfo(ctx, id)
		if err != nil {
			// Add to unresolved objects with 404 error code
			unresolvedByErrorCode[404] = append(unresolvedByErrorCode[404], id)
			bulkResponse.Summary.Unresolved++
			continue
		}

		// Check if the object requires passport authorization
		requiresAuth := false
		for _, supportedType := range obj.SupportedTypes {
			if supportedType == "Passport" {
				requiresAuth = true
				break
			}
		}

		if requiresAuth {
			// If the object requires passport authorization, add it to unresolved objects with 401 error code
			// This method doesn't provide passports for verification
			unresolvedByErrorCode[401] = append(unresolvedByErrorCode[401], id)
			bulkResponse.Summary.Unresolved++
			continue
		}

		// Add to resolved objects
		bulkResponse.ResolvedObjects = append(bulkResponse.ResolvedObjects, obj)
		bulkResponse.Summary.Resolved++
	}

	// Convert unresolved map to array
	for errorCode, objectIDs := range unresolvedByErrorCode {
		bulkResponse.UnresolvedObjects = append(bulkResponse.UnresolvedObjects, struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		}{
			ErrorCode: errorCode,
			ObjectIDs: objectIDs,
		})
	}

	// No retry needed for this operation
	return bulkResponse, nil, nil
}

// GetBulkAccessURLs retrieves multiple access URLs for objects
// This method implements the POST /access endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetBulkAccessURLs(ctx context.Context, objectAccessIDs map[string][]string) ([]models.BulkAccessURL, map[string]error) {
	// Calculate the total number of requests
	totalRequests := 0
	for _, accessIDs := range objectAccessIDs {
		totalRequests += len(accessIDs)
	}

	if totalRequests > s.config.MaxBulkRequestLength {
		return nil, map[string]error{
			"error": models.NewBadRequestError(fmt.Sprintf("Bulk request exceeds maximum limit %d", s.config.MaxBulkRequestLength)),
		}
	}

	results := []models.BulkAccessURL{}
	errors := make(map[string]error)

	for objectID, accessIDs := range objectAccessIDs {
		// Check if the object exists
		_, err := s.GetObject(ctx, objectID, nil)
		if err != nil {
			errors[objectID] = err
			continue
		}

		for _, accessID := range accessIDs {
			accessURL, err := s.repository.GetAccessURL(ctx, objectID, accessID)
			if err != nil {
				errors[fmt.Sprintf("%s:%s", objectID, accessID)] = models.NewNotFoundError(fmt.Sprintf("Access ID %s does not exist", accessID))
				continue
			}

			results = append(results, models.BulkAccessURL{
				DrsObjectID: objectID,
				DrsAccessID: accessID,
				URL:         accessURL.URL,
				Headers:     accessURL.Headers,
			})
		}
	}

	return results, errors
}

// GetBulkAccessURLsWithPassport retrieves multiple access URLs for objects with passport authorization
// This method implements the POST /access endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetBulkAccessURLsWithPassport(ctx context.Context, objectAccessIDs map[string][]string, passports []string) (*response.BulkAccessURLWithPassportResponse, error) {
	if len(objectAccessIDs) > s.config.MaxBulkRequestLength {
		return nil, models.NewBadRequestError(fmt.Sprintf("Bulk request exceeds maximum limit %d", s.config.MaxBulkRequestLength))
	}

	resp := &response.BulkAccessURLWithPassportResponse{}

	// Calculate the total number of requests
	totalRequested := 0
	for _, accessIDs := range objectAccessIDs {
		totalRequested += len(accessIDs)
	}
	resp.Summary.Requested = totalRequested

	// Initialize results
	resp.ResolvedDrsObjectAccessURLs = make([]models.BulkAccessURL, 0)

	// Error mapping for error code and object IDs
	errorMap := make(map[int][]string)

	// Global retry time
	var globalRetryAfter *int

	// Process each object and access ID
	for objectID, accessIDs := range objectAccessIDs {
		// Get object basic information
		obj, err := s.repository.GetObjectBasicInfo(ctx, objectID)
		if err != nil {
			// Object does not exist
			errorMap[404] = append(errorMap[404], objectID)
			continue
		}

		// Check if the object meets the conditions:
		// 1. SupportedTypes only contains "None"
		// 2. Has passport and can match

		// Check SupportedTypes
		hasNoneType := false
		for _, supportedType := range obj.SupportedTypes {
			if supportedType == "None" {
				hasNoneType = true
				break
			}
		}

		// Check passport authorization
		requiresAuth := len(obj.PassportAuthIssuers) > 0
		authorized := !requiresAuth // If no authorization is required, default to authorized

		if requiresAuth && len(passports) > 0 {
			validator := passport.NewValidator(obj.PassportAuthIssuers)
			for _, p := range passports {
				valid, _ := validator.ValidatePassport(p)
				if valid {
					authorized = true
					break
				}
			}
		}

		// If it's not "None" type and not authorized, skip
		if !hasNoneType && !authorized {
			errorMap[401] = append(errorMap[401], objectID)
			continue
		}

		// Process each access ID
		for _, accessID := range accessIDs {
			// Get access URL
			accessURL, retryAfter, err := s.GetAccessURLWithPassport(ctx, objectID, accessID, passports)
			if err != nil {
				// Get error code
				var statusCode int
				if drsErr, ok := err.(*models.Error); ok {
					statusCode = drsErr.StatusCode
				} else {
					statusCode = 500 // default internal server error
				}

				// Add to error group
				errorKey := fmt.Sprintf("%s:%s", objectID, accessID)
				errorMap[statusCode] = append(errorMap[statusCode], errorKey)
				continue
			}

			// Check if retry is needed
			if retryAfter != nil {
				// If retry is needed, update global retry time
				if globalRetryAfter == nil || *retryAfter < *globalRetryAfter {
					globalRetryAfter = retryAfter
				}

				// Add to unresolved list
				errorKey := fmt.Sprintf("%s:%s", objectID, accessID)
				errorMap[202] = append(errorMap[202], errorKey)
				continue
			}

			// Add to resolved list
			resp.ResolvedDrsObjectAccessURLs = append(resp.ResolvedDrsObjectAccessURLs, models.BulkAccessURL{
				DrsObjectID: objectID,
				DrsAccessID: accessID,
				URL:         accessURL.URL,
				Headers:     accessURL.Headers,
			})
		}
	}

	// Fill in unresolved object information
	for errCode, objIDs := range errorMap {
		resp.UnresolvedDrsObjects = append(resp.UnresolvedDrsObjects, struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		}{
			ErrorCode: errCode,
			ObjectIDs: objIDs,
		})
	}

	// Update summary information
	resp.Summary.Resolved = len(resp.ResolvedDrsObjectAccessURLs)
	resp.Summary.Unresolved = resp.Summary.Requested - resp.Summary.Resolved

	// Set global retry time
	if globalRetryAfter != nil {
		resp.RetryAfter = globalRetryAfter
	}

	return resp, nil
}

// GetBulkObjectsWithPassport gets multiple object information with passport authorization
// This method implements the POST /objects endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetBulkObjectsWithPassport(ctx context.Context, ids []string, expand bool, passports []string) (*response.BulkObjectResponse, *int, error) {
	// Create response structure
	bulkResponse := &response.BulkObjectResponse{
		Summary: struct {
			Requested  int `json:"requested"`
			Resolved   int `json:"resolved"`
			Unresolved int `json:"unresolved"`
		}{
			Requested:  len(ids),
			Resolved:   0,
			Unresolved: 0,
		},
		ResolvedObjects: make([]*models.DrsObject, 0),
		UnresolvedObjects: make([]struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		}, 0),
	}

	// Group unresolved objects by error code
	unresolvedByErrorCode := make(map[int][]string)

	// Process each object ID
	for _, id := range ids {
		// Get object with passport
		obj, err := s.GetObjectWithPassport(ctx, id, expand, passports)
		if err != nil {
			// Handle error
			var errorCode int
			switch e := err.(type) {
			case *models.Error:
				errorCode = e.StatusCode
			default:
				errorCode = 500
			}

			// Add to unresolved objects
			unresolvedByErrorCode[errorCode] = append(unresolvedByErrorCode[errorCode], id)
			bulkResponse.Summary.Unresolved++
		} else {
			// Add to resolved objects
			bulkResponse.ResolvedObjects = append(bulkResponse.ResolvedObjects, obj)
			bulkResponse.Summary.Resolved++
		}
	}

	// Convert unresolved map to array
	for errorCode, objectIDs := range unresolvedByErrorCode {
		bulkResponse.UnresolvedObjects = append(bulkResponse.UnresolvedObjects, struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		}{
			ErrorCode: errorCode,
			ObjectIDs: objectIDs,
		})
	}

	// No retry needed for this operation
	return bulkResponse, nil, nil
}

// ListObjects lists objects
// This method implements the GET /objects endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) ListObjects(ctx context.Context, pageSize int, pageToken string) ([]models.DrsObject, string, error) {
	// This needs to implement pagination logic, but the current repository interface does not provide a method to list all objects
	// This is a simplified implementation, and in actual applications, the repository interface should be extended
	return []models.DrsObject{}, "", models.NewNotImplementedError("List objects feature is not implemented")
}

// GetObjectAuthorization gets object authorization information
// This method implements the GET /objects/{object_id}/authz endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetObjectAuthorization(ctx context.Context, id string, passports []string) (*models.DrsObject, error) {
	// Use GetObjectAuthInfo to get authorization-related information
	obj, err := s.repository.GetObjectAuthInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	// Directly return the authorization information provided by the repository layer
	return obj, nil
}

// GetBulkObjectAuthorizations gets multiple object authorization information
// This method implements the POST /objects/authz endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetBulkObjectAuthorizations(ctx context.Context, ids []string, passports []string) (*response.BulkObjectResponse, *int, error) {
	// Create response structure
	bulkResponse := &response.BulkObjectResponse{
		Summary: struct {
			Requested  int `json:"requested"`
			Resolved   int `json:"resolved"`
			Unresolved int `json:"unresolved"`
		}{
			Requested:  len(ids),
			Resolved:   0,
			Unresolved: 0,
		},
		ResolvedObjects: make([]*models.DrsObject, 0),
		UnresolvedObjects: make([]struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		}, 0),
	}

	// Group unresolved objects by error code
	unresolvedByErrorCode := make(map[int][]string)

	// Process each object ID
	for _, id := range ids {
		// Get object authorization information
		obj, err := s.GetObjectAuthorization(ctx, id, passports)
		if err != nil {
			// Handle error
			var errorCode int
			switch e := err.(type) {
			case *models.Error:
				errorCode = e.StatusCode
			default:
				errorCode = 500
			}

			// Add to unresolved objects
			unresolvedByErrorCode[errorCode] = append(unresolvedByErrorCode[errorCode], id)
			bulkResponse.Summary.Unresolved++
		} else {
			// Add to resolved objects
			bulkResponse.ResolvedObjects = append(bulkResponse.ResolvedObjects, obj)
			bulkResponse.Summary.Resolved++
		}
	}

	// Convert unresolved map to array
	for errorCode, objectIDs := range unresolvedByErrorCode {
		bulkResponse.UnresolvedObjects = append(bulkResponse.UnresolvedObjects, struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		}{
			ErrorCode: errorCode,
			ObjectIDs: objectIDs,
		})
	}

	// No retry needed for this operation
	return bulkResponse, nil, nil
}

// GetObjectWithPassport gets object information with passport authorization
// This method implements the GET /objects/{object_id} endpoint from the GA4GH DRS API specification
func (s *DefaultDrsService) GetObjectWithPassport(ctx context.Context, id string, expand bool, passports []string) (*models.DrsObject, error) {
	// First, get object basic information
	obj, err := s.repository.GetObjectBasicInfo(ctx, id)
	if err != nil {
		return nil, models.NewNotFoundError("Object does not exist")
	}

	// Check if the object requires passport authorization
	requiresAuth := false
	for _, supportedType := range obj.SupportedTypes {
		if supportedType == "Passport" {
			requiresAuth = true
			break
		}
	}

	// If object requires passport authorization but no passport is provided, return error
	if requiresAuth && (len(passports) == 0) {
		return nil, models.NewUnauthorizedError("Passport authorization is required")
	}

	// If passport is required, validate it
	if requiresAuth {
		validator := passport.NewValidator(obj.PassportAuthIssuers)
		authorized := false

		for _, p := range passports {
			valid, _ := validator.ValidatePassport(p)
			if valid {
				authorized = true
				break
			}
		}

		if !authorized {
			return nil, models.NewUnauthorizedError("Invalid passport")
		}
	}

	// If the object information needs to be expanded
	if expand {
		// Get more detailed information, such as content lists, etc.
		// This is a sample implementation and should be customized according to actual needs
		obj.Content = append(obj.Content, models.Content{
			Type: "expanded_info",
			Data: "Additional data available",
		})
	}

	return obj, nil
}
