package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Businge931/sba-crud-ops/internal/bootstrap"
)

func main() {
	// Load configuration
	cfg := bootstrap.LoadConfig()

	// Initialize database
	pool := bootstrap.SetupDatabase(cfg)
	defer pool.Close()

	// Setup application components
	components := bootstrap.SetupComponents(cfg, pool)

	// Setup server
	server := bootstrap.SetupServer(cfg, components)

	// Try to register with API Gateway in the background
	server.RegisterWithGateway()

	// Start the server
	if err := server.Start(); err != nil {
		log.Panicf("Failed to serve: %v", err)
	}
}
