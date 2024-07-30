package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	resp, err := DrsServiceHandler(context.Background(), events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{RequestID: "handler-test"}})
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "{'response':'hello'}", resp.Body)
	}
}
