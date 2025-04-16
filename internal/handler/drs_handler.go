// internal/handler/drs_handler.go
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pennsieve/drs-service/internal/models"
	"github.com/pennsieve/drs-service/internal/models/request"
	"github.com/pennsieve/drs-service/internal/service"
)

// DrsHandler handles DRS API requests
// This struct implements HTTP handlers for all DRS API endpoints
type DrsHandler struct {
	service service.DrsService // Service layer for business logic
}

// NewDrsHandler creates a new handler instance
// It initializes the handler with the provided service
func NewDrsHandler(svc service.DrsService) *DrsHandler {
	return &DrsHandler{
		service: svc,
	}
}

// GetServiceInfo handles service information requests
// This method implements the GET /service-info endpoint from the GA4GH DRS API specification
func (h *DrsHandler) GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	info, err := h.service.GetServiceInfo(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, info)
}

// GetObject handles object retrieval requests
// This method implements the GET /objects/{object_id} endpoint from the GA4GH DRS API specification
func (h *DrsHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	objectID := vars["object_id"]

	// Get passport information from request header (if any)
	var passports []string
	if passport := r.Header.Get("GA4GH-Passport"); passport != "" {
		passports = append(passports, passport)
	}

	obj, err := h.service.GetObject(ctx, objectID, passports)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	writeJSON(w, http.StatusOK, obj)
}

// GetAccessURL handles access URL retrieval requests
// This method implements the GET /objects/{object_id}/access/{access_id} endpoint from the GA4GH DRS API specification
func (h *DrsHandler) GetAccessURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	objectID := vars["object_id"]
	accessID := vars["access_id"]

	accessURL, retryAfter, err := h.service.GetAccessURL(ctx, objectID, accessID)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	// Set status code
	statusCode := http.StatusOK
	if retryAfter != nil {
		statusCode = http.StatusAccepted
		w.Header().Set("Retry-After", fmt.Sprintf("%d", *retryAfter))
	}

	writeJSON(w, statusCode, accessURL)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	errorObj := models.Error{
		Msg:        err.Error(),
		StatusCode: status,
	}
	writeJSON(w, status, errorObj)
}

// OptionsObject handles OPTIONS requests and provides authorization information
// This method implements the OPTIONS /objects/{object_id} endpoint from the GA4GH DRS API specification
func (h *DrsHandler) OptionsObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	objectID := vars["object_id"]

	// Get passport information from request header (if any)
	var passports []string
	if passport := r.Header.Get("GA4GH-Passport"); passport != "" {
		passports = append(passports, passport)
	}

	// Use service layer method to get authorization information
	authInfo, err := h.service.GetObjectAuthorization(ctx, objectID, passports)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	// Directly use the authorization information returned by the service layer
	writeJSON(w, http.StatusOK, authInfo)
}

// PostObject handles POST requests with passport
// This method implements the POST /objects/{object_id} endpoint from the GA4GH DRS API specification
func (h *DrsHandler) PostObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	objectID := vars["object_id"]

	// Parse request body
	var req struct {
		Expand    bool     `json:"expand"`
		Passports []string `json:"passports"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Use service layer method to get object information while validating passport
	obj, err := h.service.GetObjectWithPassport(ctx, objectID, req.Expand, req.Passports)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	writeJSON(w, http.StatusOK, obj)
}

// OptionsBulkObject handles bulk OPTIONS requests
func (h *DrsHandler) OptionsBulkObject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BulkObjectIDs []string `json:"bulk_object_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	ctx := r.Context()

	// Get passport information from request header (if any)
	var passports []string
	if passport := r.Header.Get("GA4GH-Passport"); passport != "" {
		passports = append(passports, passport)
	}

	// Use service layer method to get bulk authorization information
	bulkResponse, retryAfter, err := h.service.GetBulkObjectAuthorizations(ctx, req.BulkObjectIDs, passports)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	// Set status code
	statusCode := http.StatusOK
	if bulkResponse.Summary.Unresolved > 0 {
		statusCode = http.StatusAccepted
		if retryAfter != nil {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", *retryAfter))
		}
	}

	// Directly return the response generated by the service layer
	writeJSON(w, statusCode, bulkResponse)
}

// GetBulkObjects handles bulk object retrieval requests
// This method implements the POST /objects endpoint from the GA4GH DRS API specification
func (h *DrsHandler) GetBulkObjects(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Passports     []string `json:"passports"`
		BulkObjectIDs []string `json:"bulk_object_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	ctx := r.Context()

	// Check if the request exceeds the maximum allowed size
	maxBulkSize := 100 // Should be obtained from configuration or service information
	if len(req.BulkObjectIDs) > maxBulkSize {
		writeError(w, http.StatusRequestEntityTooLarge, fmt.Errorf("bulk request exceeds maximum size of %d", maxBulkSize))
		return
	}

	// Combine passports from request body and header
	passports := req.Passports
	if passport := r.Header.Get("GA4GH-Passport"); passport != "" {
		passports = append(passports, passport)
	}

	// Process objects and collect results
	resolvedObjects := []models.DrsObject{}
	notFoundIDs := []string{}

	for _, id := range req.BulkObjectIDs {
		obj, err := h.service.GetObject(ctx, id, passports)
		if err != nil {
			notFoundIDs = append(notFoundIDs, id)
			continue
		}
		resolvedObjects = append(resolvedObjects, *obj)
	}

	// Prepare response
	response := struct {
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
	}{
		Summary: struct {
			Requested  int `json:"requested"`
			Resolved   int `json:"resolved"`
			Unresolved int `json:"unresolved"`
		}{
			Requested:  len(req.BulkObjectIDs),
			Resolved:   len(resolvedObjects),
			Unresolved: len(notFoundIDs),
		},
		ResolvedDrsObject: resolvedObjects,
	}

	// Add unresolved objects if any
	if len(notFoundIDs) > 0 {
		response.UnresolvedDrsObjects = []struct {
			ErrorCode int      `json:"error_code"`
			ObjectIDs []string `json:"object_ids"`
		}{
			{
				ErrorCode: http.StatusNotFound,
				ObjectIDs: notFoundIDs,
			},
		}
	}

	writeJSON(w, http.StatusOK, response)
}

// PostAccessURL handles POST requests for access URLs with passport
// This method implements the POST /objects/{object_id}/access/{access_id} endpoint from the GA4GH DRS API specification
func (h *DrsHandler) PostAccessURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	objectID := vars["object_id"]
	accessID := vars["access_id"]

	var req request.PostAccessURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	accessURL, retryAfter, err := h.service.GetAccessURLWithPassport(ctx, objectID, accessID, req.Passports)
	if err != nil {
		// Set status code based on error type
		if _, ok := err.(*models.Error); ok {
			writeError(w, err.(*models.Error).StatusCode, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Set status code
	statusCode := http.StatusOK
	if retryAfter != nil {
		statusCode = http.StatusAccepted
		w.Header().Set("Retry-After", fmt.Sprintf("%d", *retryAfter))
	}

	writeJSON(w, statusCode, accessURL)
}

// PostBulkAccessURL handles bulk access URL retrieval requests
func (h *DrsHandler) PostBulkAccessURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.PostBulkAccessURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Convert request format
	objectAccessIDs := make(map[string][]string)
	for _, item := range req.BulkObjectAccessIDs {
		objectAccessIDs[item.ObjectID] = item.AccessIDs
	}

	// Call service layer method
	resp, err := h.service.GetBulkAccessURLsWithPassport(ctx, objectAccessIDs, req.Passports)
	if err != nil {
		// Set status code based on error type
		if drsErr, ok := err.(*models.Error); ok {
			writeError(w, drsErr.StatusCode, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Set status code
	statusCode := http.StatusOK
	if resp.RetryAfter != nil {
		statusCode = http.StatusAccepted
		w.Header().Set("Retry-After", fmt.Sprintf("%d", *resp.RetryAfter))
	}

	writeJSON(w, statusCode, resp)
}
