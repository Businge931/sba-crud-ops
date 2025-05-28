package bootstrap

import (
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	grpcServer "github.com/Businge931/sba-crud-ops/internal/app/grpc"
	"github.com/Businge931/sba-crud-ops/proto"
)

// Server represents the gRPC server
type Server struct {
	grpcServer *grpc.Server
	config     *Config
	listener   net.Listener
}

// SetupServer initializes and configures the gRPC server
func SetupServer(cfg *Config, components *ApplicationComponents) *Server {
	// Initialize gRPC server
	server := grpc.NewServer()

	// Register odds service with improved implementation
	oddsServer := grpcServer.NewImprovedOddsServer(
		components.OddsService,
	)
	proto.RegisterOddsServiceServer(server, oddsServer)

	// Add health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus(cfg.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	// Add reflection service
	reflection.Register(server)

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ServicePort))
	if err != nil {
		log.Panicf("Failed to listen: %v", err)
	}

	return &Server{
		grpcServer: server,
		config:     cfg,
		listener:   lis,
	}
}

// Start begins serving requests
func (s *Server) Start() error {
	log.Printf("Starting %s gRPC server on port %s", s.config.ServiceName, s.config.ServicePort)
	return s.grpcServer.Serve(s.listener)
}

// RegisterWithGateway attempts to register with the API gateway
func (s *Server) RegisterWithGateway() {
	gatewayAddr := s.config.GatewayAddr
	servicePort := s.config.ServicePort

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
