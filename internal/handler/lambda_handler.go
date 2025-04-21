// internal/handler/lambda_handler.go
package handler

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/pennsieve/drs-service/internal/config"
	"github.com/pennsieve/drs-service/internal/logging"
	"github.com/pennsieve/drs-service/internal/repository"
	"github.com/pennsieve/drs-service/internal/router"
	"github.com/pennsieve/drs-service/internal/service"
)

var logger = logging.Default
var serviceConfig *config.Config
var drsService service.DrsService
var drsHandler *DrsHandler

func init() {
	// Load configuration
	serviceConfig = config.LoadConfig()

	// Initialize repository
	repo, err := repository.NewMockDrsRepository()
	if err != nil {
		logger.Error("Failed to initialize repository", "error", err)
		return
	}

	// Initialize service layer
	drsService = service.NewDrsService(repo, serviceConfig.BaseURL, serviceConfig)

	// Initialize handler
	drsHandler = NewDrsHandler(drsService)
}

// DrsServiceHandler is the main entry point for AWS Lambda
func DrsServiceHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Add request ID to logger
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		logger = logger.With(
			slog.String("requestID", lc.AwsRequestID),
			slog.String("handler", "DrsServiceHandler"),
		)
	}

	logger.Info("request received",
		"routeKey", request.RouteKey,
		"pathParameters", request.PathParameters,
		"rawPath", request.RawPath)

	// Create router
	router := router.NewLambdaRouter(logger)

	// GA4GH DRS API endpoints
	router.GET("/ga4gh/drs/v1/service-info", GetServiceInfoHandler)
	router.GET("/ga4gh/drs/v1/objects/{object_id}", GetObjectHandler)
	router.POST("/ga4gh/drs/v1/objects/{object_id}", PostObjectHandler)
	router.OPTIONS("/ga4gh/drs/v1/objects/{object_id}", OptionsObjectHandler)

	router.POST("/ga4gh/drs/v1/objects", GetBulkObjectsHandler)
	router.OPTIONS("/ga4gh/drs/v1/objects", OptionsBulkObjectHandler)

	router.GET("/ga4gh/drs/v1/objects/{object_id}/access/{access_id}", GetAccessURLHandler)
	router.POST("/ga4gh/drs/v1/objects/{object_id}/access/{access_id}", PostAccessURLHandler)

	router.POST("/ga4gh/drs/v1/objects/access", PostBulkAccessURLHandler)

	return router.Start(ctx, request)
}
