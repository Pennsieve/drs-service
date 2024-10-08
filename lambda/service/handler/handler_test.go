package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pennsieve/drs-service/service/models"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("DRS_SERVICE_ID", "io.pennsieve.drs")
	os.Setenv("DRS_ORG_URL", "https://pennsieve.io")

	code := m.Run()
	os.Exit(code)
}

func TestDrsServiceHandler_ServiceInfo(t *testing.T) {
	request := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/ga4gh/drs/v1/service-info",
			},
		},
	}

	resp, err := DrsServiceHandler(context.Background(), request)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var serviceInfo models.ServiceInfo
		err = json.Unmarshal([]byte(resp.Body), &serviceInfo)
		assert.NoError(t, err)

		expectedServiceInfo := models.NewServiceInfo()
		assert.Equal(t, expectedServiceInfo, serviceInfo)
	}
}

func TestDrsServiceHandler_NotFound(t *testing.T) {

	request := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path: "/unknown-path",
			},
		},
	}

	resp, err := DrsServiceHandler(context.Background(), request)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	assert.Equal(t, `{"error":"Route not found"}`, resp.Body)
}
