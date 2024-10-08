package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/pennsieve/drs-service/service/models"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func DrsServiceHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		logger = logger.With(
			slog.String("requestID", lc.AwsRequestID),
			slog.String("handler", "DrsServiceHandler"),
		)
	}

	router := NewLambdaRouter()
	router.GET("/ga4gh/drs/v1/service-info", handleServiceInfoRequest)
	return router.Start(ctx, request)
}

func handleServiceInfoRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	logger.Info("handleServiceInfoRequest()")
	serviceInfo := models.NewServiceInfo()
	body, err := json.Marshal(serviceInfo)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"msg": "Internal Server Error", "status_code": 500}`,
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}
