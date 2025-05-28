package bootstrap

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/Businge931/sba-crud-ops/internal/core/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type mockOddsService struct {
	ports.OddsService
}

func TestServerSetup(t *testing.T) {
	cfg := &Config{
		ServicePort: "0", // Use port 0 to get a random available port
		ServiceName: "test-service",
	}

	components := &ApplicationComponents{
		OddsService:   &mockOddsService{},
		OddsValidator: validator.NewDefaultOddsValidator(NewLeagueRegistry(cfg)),
	}

	// Test server setup
	t.Run("server setup and start", func(t *testing.T) {
		server := SetupServer(cfg, components)
		require.NotNil(t, server)
		require.NotNil(t, server.grpcServer)

		// Start server in a goroutine
		errChan := make(chan error, 1)
		go func() {
			errChan <- server.Start()
		}()

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Test health check
		testHealthCheck(t, server.listener.Addr().String())

		// Graceful shutdown
		server.grpcServer.GracefulStop()

		// Verify server stopped
		select {
		case err := <-errChan:
			assert.NoError(t, err)
		case <-time.After(1 * time.Second):
			t.Fatal("Server did not stop within timeout")
		}
	})
}

func testHealthCheck(t *testing.T, addr string) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	resp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	require.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
}

func TestServerRegistration(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "successful registration",
			config: &Config{
				ServiceName: "test-service",
				GatewayAddr: "localhost:8080",
			},
			expectError: false,
		},
		{
			name: "invalid gateway address",
			config: &Config{
				ServiceName: "test-service",
				GatewayAddr: "invalid-address",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{
				config:     tt.config,
				grpcServer: grpc.NewServer(),
			}

			// Setup a listener
			lis, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			server.listener = lis

			// Note: RegisterWithGateway is not directly testable in unit tests
			// as it requires a running gateway server. This is better suited for
			// integration tests.
			t.Skip("Skipping RegisterWithGateway test as it requires a running gateway server")
		})
	}
}

func TestServerStartError(t *testing.T) {
	// Create a server with an invalid port
	server := &Server{
		config:     &Config{ServicePort: "99999"}, // Invalid port number
		grpcServer: grpc.NewServer(),
	}

	// Test that Start panics with an invalid port
	assert.Panics(t, func() {
		err := server.Start()
		if err != nil {
			panic(err)
		}
	}, "Expected Start to panic with invalid port")
}
