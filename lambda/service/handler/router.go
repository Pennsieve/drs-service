package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type RouterHandlerFunc func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

type Router interface {
	POST(string, RouterHandlerFunc)
	GET(string, RouterHandlerFunc)
	DELETE(string, RouterHandlerFunc)
	PUT(string, RouterHandlerFunc)
	Start(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
}

type LambdaRouter struct {
	getRoutes    map[string]RouterHandlerFunc
	postRoutes   map[string]RouterHandlerFunc
	deleteRoutes map[string]RouterHandlerFunc
	putRoutes    map[string]RouterHandlerFunc
}

func NewLambdaRouter() Router {
	return &LambdaRouter{
		getRoutes:    make(map[string]RouterHandlerFunc),
		postRoutes:   make(map[string]RouterHandlerFunc),
		deleteRoutes: make(map[string]RouterHandlerFunc),
		putRoutes:    make(map[string]RouterHandlerFunc),
	}
}

func (r *LambdaRouter) POST(routeKey string, handler RouterHandlerFunc) {
	r.postRoutes[routeKey] = handler
}

func (r *LambdaRouter) GET(routeKey string, handler RouterHandlerFunc) {
	r.getRoutes[routeKey] = handler
}

func (r *LambdaRouter) DELETE(routeKey string, handler RouterHandlerFunc) {
	r.deleteRoutes[routeKey] = handler
}

func (r *LambdaRouter) PUT(routeKey string, handler RouterHandlerFunc) {
	r.putRoutes[routeKey] = handler
}

func (r *LambdaRouter) Start(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	routeKey := request.RequestContext.HTTP.Path
	logger.Info("Routing request", "method", request.RequestContext.HTTP.Method, "path", routeKey)

	switch request.RequestContext.HTTP.Method {
	case http.MethodPost:
		f, ok := r.postRoutes[routeKey]
		if ok {
			return f(ctx, request)
		}
	case http.MethodGet:
		f, ok := r.getRoutes[routeKey]
		if ok {
			return f(ctx, request)
		}
	case http.MethodDelete:
		f, ok := r.deleteRoutes[routeKey]
		if ok {
			return f(ctx, request)
		}
	case http.MethodPut:
		f, ok := r.putRoutes[routeKey]
		if ok {
			return f(ctx, request)
		}
	}

	logger.Info("Route not found", "path", routeKey)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"error":"Route not found"}`,
	}, nil
}
