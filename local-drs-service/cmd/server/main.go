// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/daniellong/local-drs-service/config"
	"github.com/daniellong/local-drs-service/internal/handler"
	"github.com/daniellong/local-drs-service/internal/repository"
	"github.com/daniellong/local-drs-service/internal/router"
	"github.com/daniellong/local-drs-service/internal/service"
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
