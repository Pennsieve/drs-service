// cmd/server/main.go
/*
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pennsieve/drs-service/internal/config"
	"github.com/pennsieve/drs-service/internal/handler"
	"github.com/pennsieve/drs-service/internal/repository"
	"github.com/pennsieve/drs-service/internal/router"
	"github.com/pennsieve/drs-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize repository
	repo, err := repository.NewMockDrsRepository()
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize service layer
	drsService := service.NewDrsService(repo, cfg.BaseURL, cfg)

	// Initialize handler
	drsHandler := handler.NewDrsHandler(drsService)

	// Setup routes
	r := router.SetupRoutes(drsHandler)

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
*/
// cmd/lambda/main.go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pennsieve/drs-service/internal/handler"
)

func main() {
	lambda.Start(handler.DrsServiceHandler)
}
