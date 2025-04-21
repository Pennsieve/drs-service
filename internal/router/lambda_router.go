// internal/router/lambda_router.go
package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type RouterHandlerFunc func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

type LambdaRouter interface {
	GET(string, RouterHandlerFunc)
	POST(string, RouterHandlerFunc)
	PUT(string, RouterHandlerFunc)
	DELETE(string, RouterHandlerFunc)
	OPTIONS(string, RouterHandlerFunc)
	Start(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
}

type lambdaRouter struct {
	getRoutes     map[string]RouterHandlerFunc
	postRoutes    map[string]RouterHandlerFunc
	putRoutes     map[string]RouterHandlerFunc
	deleteRoutes  map[string]RouterHandlerFunc
	optionsRoutes map[string]RouterHandlerFunc
	logger        Logger
}

// Logger interface to allow dependency injection for testing
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// NewLambdaRouter creates a new router for AWS Lambda
func NewLambdaRouter(logger Logger) LambdaRouter {
	return &lambdaRouter{
		getRoutes:     make(map[string]RouterHandlerFunc),
		postRoutes:    make(map[string]RouterHandlerFunc),
		putRoutes:     make(map[string]RouterHandlerFunc),
		deleteRoutes:  make(map[string]RouterHandlerFunc),
		optionsRoutes: make(map[string]RouterHandlerFunc),
		logger:        logger,
	}
}

func (r *lambdaRouter) GET(path string, handler RouterHandlerFunc) {
	r.getRoutes[path] = handler
}

func (r *lambdaRouter) POST(path string, handler RouterHandlerFunc) {
	r.postRoutes[path] = handler
}

func (r *lambdaRouter) PUT(path string, handler RouterHandlerFunc) {
	r.putRoutes[path] = handler
}

func (r *lambdaRouter) DELETE(path string, handler RouterHandlerFunc) {
	r.deleteRoutes[path] = handler
}

func (r *lambdaRouter) OPTIONS(path string, handler RouterHandlerFunc) {
	r.optionsRoutes[path] = handler
}

func (r *lambdaRouter) Start(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	r.logger.Info("Router processing request",
		"method", request.RequestContext.HTTP.Method,
		"path", request.RawPath)

	path := request.RawPath

	// Find matching route handler
	var handler RouterHandlerFunc
	var found bool

	switch request.RequestContext.HTTP.Method {
	case http.MethodGet:
		handler, found = r.findHandler(r.getRoutes, path, &request)
	case http.MethodPost:
		handler, found = r.findHandler(r.postRoutes, path, &request)
	case http.MethodPut:
		handler, found = r.findHandler(r.putRoutes, path, &request)
	case http.MethodDelete:
		handler, found = r.findHandler(r.deleteRoutes, path, &request)
	case http.MethodOptions:
		handler, found = r.findHandler(r.optionsRoutes, path, &request)
	default:
		r.logger.Warn("Unsupported HTTP method", "method", request.RequestContext.HTTP.Method)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       `{"error": "Method not allowed"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	if !found {
		r.logger.Warn("Route not found", "path", path)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Route not found"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return handler(ctx, request)
}

// Find matching route handler with path parameters support
func (r *lambdaRouter) findHandler(routes map[string]RouterHandlerFunc, path string, request *events.APIGatewayV2HTTPRequest) (RouterHandlerFunc, bool) {
	// Exact match
	if handler, ok := routes[path]; ok {
		return handler, true
	}

	// Try to match routes with parameters
	for pattern, handler := range routes {
		if params, match := r.matchRoute(pattern, path); match {
			// If PathParameters is nil, initialize it
			if request.PathParameters == nil {
				request.PathParameters = make(map[string]string)
			}

			// Add extracted parameters to request
			for k, v := range params {
				request.PathParameters[k] = v
			}

			return handler, true
		}
	}

	return nil, false
}

// Route matching logic with {param} format path parameters
// Returns extracted parameters and match status
func (r *lambdaRouter) matchRoute(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i, part := range patternParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			// This is a parameter part, extract parameter name and store value
			paramName := part[1 : len(part)-1]
			params[paramName] = pathParts[i]
			continue
		}

		if part != pathParts[i] {
			return nil, false
		}
	}

	return params, true
}
