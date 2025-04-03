package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	grpcServer "github.com/Businge931/sba-crud-ops/internal/app/grpc"
	"github.com/Businge931/sba-crud-ops/internal/db"
	"github.com/Businge931/sba-crud-ops/internal/env"
	"github.com/Businge931/sba-crud-ops/internal/repository/postgres"
	"github.com/Businge931/sba-crud-ops/internal/service"
	"github.com/Businge931/sba-crud-ops/proto"

	log "github.com/sirupsen/logrus"
)

const (
	version            = "0.0.1"
	serviceName        = "odds-service"
	defaultServicePort = "50052"
)

func main() {
	// Load configuration
	servicePort := env.GetString("SERVICE_PORT", defaultServicePort)
	gatewayAddr := env.GetString("GATEWAY_ADDR", "localhost:8080")
	dbAddr := env.GetString("DB_ADDR", "postgresql://admin:adminpassword@localhost:5433/sba_crud_ops?sslmode=disable")
	maxOpenConns := env.GetInt("DB_MAX_OPEN_CONNS", 30)
	maxIdleConns := env.GetInt("DB_MAX_IDLE_CONNS", 30)
	maxIdleTime := env.GetString("DB_MAX_IDLE_TIME", "15m")

	// Initialize database connection pool
	pool, err := db.New(dbAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Database connection pool established")

	// Initialize repository
	oddsRepo := postgres.NewOddsRepository(pool)

	// Initialize service
	oddsService := service.NewOddsService(oddsRepo)

	// Initialize gRPC server
	server := grpc.NewServer()
	oddsServer := grpcServer.NewOddsServer(oddsService)
	proto.RegisterOddsServiceServer(server, oddsServer)

	// Add health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_SERVING)

	// Add reflection service
	reflection.Register(server)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", servicePort))
	if err != nil {
		log.Panicf("Failed to listen: %v", err)
	}

	log.Printf("Starting %s gRPC server on port %s", serviceName, servicePort)

	// Try to register with API Gateway in the background
	registerWithGateway(gatewayAddr, servicePort)

	if err := server.Serve(lis); err != nil {
		log.Panicf("Failed to serve: %v", err)
	}
}

// registerWithGateway attempts to connect to the API gateway and register this service
// but does not prevent the service from starting if registration fails
func registerWithGateway(gatewayAddr, servicePort string) {
	// Don't block service startup if registration fails
	go func() {
		maxRetries := 3
		retryInterval := 5 * time.Second

		for i := 0; i < maxRetries; i++ {
			log.Printf("Attempting to register with API gateway at %s (attempt %d/%d)", gatewayAddr, i+1, maxRetries)

			// Skip registration if gateway address is the default localhost:8080
			// as it's likely running HTTP instead of gRPC on that port
			if gatewayAddr == "localhost:8080" {
				log.Printf("Default gateway address detected. Skipping automatic registration.")
				log.Printf("Service is available at port %s. Make sure the API gateway is configured to use it.", servicePort)
				return
			}

			// Connect to API Gateway
			conn, err := grpc.Dial(
				gatewayAddr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				log.Printf("Failed to connect to API Gateway: %v. Retrying in %s...", err, retryInterval)
				time.Sleep(retryInterval)
				continue
			}

			// Create client for gateway registration service
			client := proto.NewOddsServiceClient(conn)

			// Attempt to verify connectivity
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err = client.GetOdds(ctx, &proto.GetOddsRequest{OddsId: "test-connection"})
			if err != nil {
				log.Printf("API Gateway connection test failed: %v. Retrying in %s...", err, retryInterval)
				conn.Close()
				time.Sleep(retryInterval)
				continue
			}

			log.Printf("Successfully connected to API Gateway at %s", gatewayAddr)
			conn.Close()
			return
		}

		log.Printf("NOTE: Could not establish connection with API gateway after %d attempts", maxRetries)
		log.Printf("Service is running on port %s and will be available when the API gateway connects to it", servicePort)
	}()
}
