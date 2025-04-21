// internal/handler/lambda_handlers.go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pennsieve/drs-service/internal/models"
)

// Lambda handler functions that adapt the existing DrsHandler methods to Lambda format

// GetServiceInfoHandler handles service information requests
func GetServiceInfoHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "GetServiceInfoHandler"))
	handlerLogger.Info("Getting service info")

	info, err := drsService.GetServiceInfo(ctx)
	if err != nil {
		handlerLogger.Error("Failed to get service info", "error", err)
		return createErrorResponse(http.StatusInternalServerError, err)
	}

	return createJSONResponse(http.StatusOK, info)
}

// GetObjectHandler handles object retrieval requests
func GetObjectHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "GetObjectHandler"))

	objectID, ok := request.PathParameters["object_id"]
	if !ok || objectID == "" {
		handlerLogger.Error("Missing object_id parameter")
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("missing object_id parameter"))
	}

	handlerLogger.Info("Getting object", "object_id", objectID)

	// Get passport information from request header (if any)
	var passports []string
	if passport, ok := request.Headers["ga4gh-passport"]; ok && passport != "" {
		passports = append(passports, passport)
	}

	obj, err := drsService.GetObject(ctx, objectID, passports)
	if err != nil {
		handlerLogger.Error("Failed to get object", "error", err, "object_id", objectID)
		return createErrorResponse(http.StatusNotFound, err)
	}

	return createJSONResponse(http.StatusOK, obj)
}

// PostObjectHandler handles POST requests with passport
func PostObjectHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "PostObjectHandler"))

	objectID, ok := request.PathParameters["object_id"]
	if !ok || objectID == "" {
		handlerLogger.Error("Missing object_id parameter")
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("missing object_id parameter"))
	}

	// Parse request body
	var req struct {
		Expand    bool     `json:"expand"`
		Passports []string `json:"passports"`
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		handlerLogger.Error("Invalid request body", "error", err)
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
	}

	handlerLogger.Info("Getting object with passport", "object_id", objectID)

	obj, err := drsService.GetObjectWithPassport(ctx, objectID, req.Expand, req.Passports)
	if err != nil {
		handlerLogger.Error("Failed to get object with passport", "error", err, "object_id", objectID)
		return createErrorResponse(http.StatusNotFound, err)
	}

	return createJSONResponse(http.StatusOK, obj)
}

// OptionsObjectHandler handles OPTIONS requests and provides authorization information
func OptionsObjectHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "OptionsObjectHandler"))

	objectID, ok := request.PathParameters["object_id"]
	if !ok || objectID == "" {
		handlerLogger.Error("Missing object_id parameter")
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("missing object_id parameter"))
	}

	// Get passport information from request header (if any)
	var passports []string
	if passport, ok := request.Headers["ga4gh-passport"]; ok && passport != "" {
		passports = append(passports, passport)
	}

	handlerLogger.Info("Getting object authorization", "object_id", objectID)

	authInfo, err := drsService.GetObjectAuthorization(ctx, objectID, passports)
	if err != nil {
		handlerLogger.Error("Failed to get object authorization", "error", err, "object_id", objectID)
		return createErrorResponse(http.StatusNotFound, err)
	}

	return createJSONResponse(http.StatusOK, authInfo)
}

// GetBulkObjectsHandler handles bulk object retrieval requests
func GetBulkObjectsHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "GetBulkObjectsHandler"))

	var req struct {
		Passports     []string `json:"passports"`
		BulkObjectIDs []string `json:"bulk_object_ids"`
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		handlerLogger.Error("Invalid request body", "error", err)
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
	}

	handlerLogger.Info("Getting bulk objects", "count", len(req.BulkObjectIDs))

	// Note: According to error messages, GetBulkObjects only accepts two parameters, not three
	bulkResponse, retryAfter, err := drsService.GetBulkObjects(ctx, req.BulkObjectIDs)
	if err != nil {
		handlerLogger.Error("Failed to get bulk objects", "error", err)
		return createErrorResponse(http.StatusBadRequest, err)
	}

	// Set status code
	statusCode := http.StatusOK
	headers := map[string]string{"Content-Type": "application/json"}

	if bulkResponse.Summary.Unresolved > 0 {
		statusCode = http.StatusAccepted
		if retryAfter != nil {
			headers["Retry-After"] = strconv.Itoa(*retryAfter)
		}
	}

	return createJSONResponseWithHeaders(statusCode, bulkResponse, headers)
}

// OptionsBulkObjectHandler handles bulk OPTIONS requests
func OptionsBulkObjectHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "OptionsBulkObjectHandler"))

	var req struct {
		BulkObjectIDs []string `json:"bulk_object_ids"`
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		handlerLogger.Error("Invalid request body", "error", err)
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
	}

	// Get passport information from request header (if any)
	var passports []string
	if passport, ok := request.Headers["ga4gh-passport"]; ok && passport != "" {
		passports = append(passports, passport)
	}

	handlerLogger.Info("Getting bulk object authorizations", "count", len(req.BulkObjectIDs))

	bulkResponse, retryAfter, err := drsService.GetBulkObjectAuthorizations(ctx, req.BulkObjectIDs, passports)
	if err != nil {
		handlerLogger.Error("Failed to get bulk object authorizations", "error", err)
		return createErrorResponse(http.StatusBadRequest, err)
	}

	// Set status code
	statusCode := http.StatusOK
	headers := map[string]string{"Content-Type": "application/json"}

	if bulkResponse.Summary.Unresolved > 0 {
		statusCode = http.StatusAccepted
		if retryAfter != nil {
			headers["Retry-After"] = strconv.Itoa(*retryAfter)
		}
	}

	return createJSONResponseWithHeaders(statusCode, bulkResponse, headers)
}

// GetAccessURLHandler handles access URL retrieval requests
func GetAccessURLHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "GetAccessURLHandler"))

	objectID, ok := request.PathParameters["object_id"]
	if !ok || objectID == "" {
		handlerLogger.Error("Missing object_id parameter")
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("missing object_id parameter"))
	}

	accessID, ok := request.PathParameters["access_id"]
	if !ok || accessID == "" {
		handlerLogger.Error("Missing access_id parameter")
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("missing access_id parameter"))
	}

	handlerLogger.Info("Getting access URL", "object_id", objectID, "access_id", accessID)

	accessURL, retryAfter, err := drsService.GetAccessURL(ctx, objectID, accessID)
	if err != nil {
		handlerLogger.Error("Failed to get access URL", "error", err, "object_id", objectID, "access_id", accessID)
		return createErrorResponse(http.StatusNotFound, err)
	}

	// Set status code
	statusCode := http.StatusOK
	headers := map[string]string{"Content-Type": "application/json"}

	if retryAfter != nil {
		statusCode = http.StatusAccepted
		headers["Retry-After"] = strconv.Itoa(*retryAfter)
	}

	return createJSONResponseWithHeaders(statusCode, accessURL, headers)
}

// PostAccessURLHandler handles POST requests for access URLs with passport
func PostAccessURLHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "PostAccessURLHandler"))

	objectID, ok := request.PathParameters["object_id"]
	if !ok || objectID == "" {
		handlerLogger.Error("Missing object_id parameter")
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("missing object_id parameter"))
	}

	accessID, ok := request.PathParameters["access_id"]
	if !ok || accessID == "" {
		handlerLogger.Error("Missing access_id parameter")
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("missing access_id parameter"))
	}

	// Parse request body
	var req struct {
		Passports []string `json:"passports"`
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		handlerLogger.Error("Invalid request body", "error", err)
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
	}

	handlerLogger.Info("Getting access URL with passport", "object_id", objectID, "access_id", accessID)

	accessURL, retryAfter, err := drsService.GetAccessURLWithPassport(ctx, objectID, accessID, req.Passports)
	if err != nil {
		handlerLogger.Error("Failed to get access URL with passport", "error", err, "object_id", objectID, "access_id", accessID)
		return createErrorResponse(http.StatusNotFound, err)
	}

	// Set status code
	statusCode := http.StatusOK
	headers := map[string]string{"Content-Type": "application/json"}

	if retryAfter != nil {
		statusCode = http.StatusAccepted
		headers["Retry-After"] = strconv.Itoa(*retryAfter)
	}

	return createJSONResponseWithHeaders(statusCode, accessURL, headers)
}

// PostBulkAccessURLHandler handles bulk access URL retrieval requests
func PostBulkAccessURLHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	handlerLogger := logger.With(slog.String("handler", "PostBulkAccessURLHandler"))

	// Parse request body
	var req struct {
		Passports []string `json:"passports"`
		Requests  []struct {
			ObjectID string `json:"object_id"`
			AccessID string `json:"access_id"`
		} `json:"requests"`
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		handlerLogger.Error("Invalid request body", "error", err)
		return createErrorResponse(http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
	}

	handlerLogger.Info("Getting bulk access URLs", "count", len(req.Requests))

	// Convert request to map[string][]string format
	objectAccessIDs := make(map[string][]string)
	for _, r := range req.Requests {
		objectAccessIDs[r.ObjectID] = append(objectAccessIDs[r.ObjectID], r.AccessID)
	}

	// Call service method to get bulk access URLs
	bulkURLs, errMap := drsService.GetBulkAccessURLs(ctx, objectAccessIDs)

	// Check for errors
	if len(errMap) > 0 {
		// If all requests failed, return an error
		if len(errMap) == len(objectAccessIDs) {
			handlerLogger.Error("All bulk access URL requests failed", "errors", errMap)
			// Choose the first error as an example
			var firstErr error
			for _, err := range errMap {
				firstErr = err
				break
			}
			return createErrorResponse(http.StatusBadRequest, firstErr)
		}
		// Otherwise, log errors but continue processing
		handlerLogger.Warn("Some bulk access URL requests failed", "errors", errMap)
	}

	// Create response
	response := struct {
		AccessURLs []models.BulkAccessURL `json:"access_urls"`
		Errors     map[string]string      `json:"errors,omitempty"`
	}{
		AccessURLs: bulkURLs,
	}

	// If there are errors, add them to the response
	if len(errMap) > 0 {
		errorMessages := make(map[string]string)
		for objID, err := range errMap {
			errorMessages[objID] = err.Error()
		}
		response.Errors = errorMessages
	}

	// Set status code
	statusCode := http.StatusOK
	headers := map[string]string{"Content-Type": "application/json"}

	// If there are partial failures, use 202 Accepted
	if len(errMap) > 0 {
		statusCode = http.StatusAccepted
	}

	return createJSONResponseWithHeaders(statusCode, response, headers)
}

// Helper functions for creating responses

func createJSONResponse(statusCode int, data interface{}) (events.APIGatewayV2HTTPResponse, error) {
	return createJSONResponseWithHeaders(statusCode, data, map[string]string{"Content-Type": "application/json"})
}

func createJSONResponseWithHeaders(statusCode int, data interface{}, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal response", "error", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Internal server error"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers:    headers,
	}, nil
}

func createErrorResponse(statusCode int, err error) (events.APIGatewayV2HTTPResponse, error) {
	errorObj := models.Error{
		Msg:        err.Error(),
		StatusCode: statusCode,
	}

	return createJSONResponse(statusCode, errorObj)
}
