package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestDrsServiceHandler_ServiceInfo(t *testing.T) {
	resp, err := DrsServiceHandler(context.Background(), events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path: "/ga4gh/drs/v1/service-info",
			},
			RequestID: "handler-test",
		},
	})


	if assert.NoError(t, err) {

		assert.Equal(t, http.StatusOK, resp.StatusCode)


		expectedBody := `{
			"id":"com.example.drs",
			"name":"Example DRS Service",
			"type":{"group":"org.ga4gh","artifact":"drs"},
			"description":"This service provides DRS functionalities.",
			"organization":{"name":"Example Organization","url":"https://example.org"},
			"contactUrl":"mailto:support@example.org",
			"documentationUrl":"https://example.org/docs",
			"createdAt":"2024-09-30T00:00:00Z",
			"updatedAt":"2024-09-30T00:00:00Z"
		}`

		assert.JSONEq(t, expectedBody, resp.Body)
	}
}

func TestDrsServiceHandler_NotFound(t *testing.T) {

	resp, err := DrsServiceHandler(context.Background(), events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path: "/unknown-path",
			},
			RequestID: "handler-test",
		},
	})


	if assert.NoError(t, err) {

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, `{"error":"Not Found"}`, resp.Body)
	}
}
