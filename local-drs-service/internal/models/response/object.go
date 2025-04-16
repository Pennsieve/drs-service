// internal/models/response/object.go
package response

import (
	model "github.com/daniellong/local-drs-service/internal/models"
)

// ObjectResponse represents a response for object retrieval
// This struct encapsulates the DrsObject model from GA4GH DRS API v1.4.0
// Used for:
// 1. In the response of GET /objects/{object_id} endpoint
// 2. In the response of POST /objects/{object_id} endpoint with Passport
//
// Field description:
// - DrsObject: The DRS object data (required)
//   Implemented by embedding model.DrsObject to maintain consistency with the model layer
type ObjectResponse struct {
	*model.DrsObject
}

// CreateObjectResponse represents a response for object creation
// This struct encapsulates the DrsObject model for creation operations
// Used for:
// 1. In the response of object creation operations (internal use)
//
// Field description:
// - DrsObject: The created DRS object data (required)
//   Implemented by embedding model.DrsObject to maintain consistency with the model layer
type CreateObjectResponse struct {
	*model.DrsObject
}
