// internal/router/router.go
package router

import (
	"github.com/gorilla/mux"

	"github.com/daniellong/local-drs-service/internal/handler"
)

// SetupRoutes configures all routes for the DRS API
// It registers all endpoints defined in the GA4GH DRS API specification
func SetupRoutes(drsHandler *handler.DrsHandler) *mux.Router {
	// Create main router
	r := mux.NewRouter()

	// Create API subrouter with the required prefix
	apiRouter := r.PathPrefix("/ga4gh/drs/v1").Subrouter()

	// Service info endpoint
	apiRouter.HandleFunc("/service-info", drsHandler.GetServiceInfo).Methods("GET")

	// Object endpoints
	apiRouter.HandleFunc("/objects/{object_id}", drsHandler.GetObject).Methods("GET")
	apiRouter.HandleFunc("/objects/{object_id}", drsHandler.PostObject).Methods("POST")
	apiRouter.HandleFunc("/objects/{object_id}", drsHandler.OptionsObject).Methods("OPTIONS")

	// Bulk object endpoints
	apiRouter.HandleFunc("/objects", drsHandler.GetBulkObjects).Methods("POST")
	apiRouter.HandleFunc("/objects", drsHandler.OptionsBulkObject).Methods("OPTIONS")

	// Access URL endpoints
	apiRouter.HandleFunc("/objects/{object_id}/access/{access_id}", drsHandler.GetAccessURL).Methods("GET")
	apiRouter.HandleFunc("/objects/{object_id}/access/{access_id}", drsHandler.PostAccessURL).Methods("POST")

	// Passport-based bulk access URL endpoint
	apiRouter.HandleFunc("/objects/access", drsHandler.PostBulkAccessURL).Methods("POST")

	return r
}
