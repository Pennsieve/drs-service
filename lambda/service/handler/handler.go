package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pennsieve/drs-service/service/logging"
)


type ServiceInfo struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Type           TypeInfo `json:"type"`
	Description    string   `json:"description"`
	Organization   OrgInfo  `json:"organization"`
	ContactURL     string   `json:"contactUrl"`
	DocumentationURL string `json:"documentationUrl"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type TypeInfo struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
}

type OrgInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}


var drsServiceInfo = ServiceInfo{
	ID:          "com.example.drs",
	Name:        "Example DRS Service",
	Type:        TypeInfo{Group: "org.ga4gh", Artifact: "drs"},
	Description: "This service provides DRS functionalities.",
	Organization: OrgInfo{
		Name: "Example Organization",
		URL:  "https://example.org",
	},
	ContactURL:     "mailto:support@example.org",
	DocumentationURL: "https://example.org/docs",
	CreatedAt:      "2024-09-30T00:00:00Z",
	UpdatedAt:      "2024-09-30T00:00:00Z",
}

var logger = logging.Default

func init() {
	logger.Info("init()")
}


func DrsServiceHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	logger = logger.With(slog.String("requestID", request.RequestContext.RequestID))


	if request.RequestContext.HTTP.Path == "/ga4gh/drs/v1/service-info" {
		return handleServiceInfoRequest()
	}


	apiResponse := events.APIGatewayV2HTTPResponse{
		Body:       `{"error":"Not Found"}`,
		StatusCode: 404,
	}
	return &apiResponse, nil
}


func handleServiceInfoRequest() (*events.APIGatewayV2HTTPResponse, error) {
	logger.Info("handleServiceInfoRequest()")


	body, err := json.Marshal(drsServiceInfo)
	if err != nil {

		apiResponse := events.APIGatewayV2HTTPResponse{
			Body:       `{"msg": "Internal Server Error", "status_code": 500}`,
			StatusCode: 500,
		}
		return &apiResponse, err 
	}


	apiResponse := events.APIGatewayV2HTTPResponse{
		Body:       string(body),
		StatusCode: 200,
	}
	return &apiResponse, nil
}
