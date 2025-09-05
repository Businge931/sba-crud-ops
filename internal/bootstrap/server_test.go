package bootstrap

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/adoptors/secondary/validator"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/Businge931/sba-crud-ops/proto"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type serverTestDeps struct {
	server  *Server
	errChan chan error
}

type serverTestArgs struct {
	config *Config
}

type serverTest struct {
	name     string
	deps     serverTestDeps
	args     serverTestArgs
	before   func(*testing.T, *serverTestDeps, *serverTestArgs)
	after    func(*testing.T, *serverTestDeps, *serverTestArgs)
	validate func(*testing.T, *serverTestDeps, *serverTestArgs, error)
}

func TestServerSetup(t *testing.T) {

	tests := []serverTest{
		{
			name: "successful server setup",
			deps: serverTestDeps{},
			args: serverTestArgs{
				config: &Config{
					ServiceName: "test-service",
					ServicePort: "0", // Let OS choose an available port
				},
			},
			before: func(t *testing.T, deps *serverTestDeps, args *serverTestArgs) {
				deps.errChan = make(chan error, 1)
				deps.server = SetupServer(args.config, &ApplicationComponents{
					OddsService:   &mockOddsService{},
					OddsValidator: validator.NewDefaultOddsValidator(NewLeagueRegistry(args.config)),
				})

				go func() {
					deps.errChan <- deps.server.Start()
				}()
				time.Sleep(100 * time.Millisecond) // Give server time to start
			},
			after: func(t *testing.T, deps *serverTestDeps, args *serverTestArgs) {
				if deps.server != nil && deps.server.grpcServer != nil {
					deps.server.grpcServer.GracefulStop()
					select {
					case err := <-deps.errChan:
						assert.NoError(t, err, "Server returned an error during graceful stop")
					case <-time.After(1 * time.Second):
						t.Error("Server did not stop within timeout")
					}
				}
			},
			validate: func(t *testing.T, deps *serverTestDeps, args *serverTestArgs, err error) {
				assert.NoError(t, err, "Server setup should not return an error")
				assert.NotNil(t, deps.server, "Server should not be nil")
				assert.NotNil(t, deps.server.listener, "Listener should be initialized")
				testHealthCheck(t, deps.server.listener.Addr().String())
			},
		},
		{
			name: "fails with invalid port",
			deps: serverTestDeps{},
			args: serverTestArgs{
				config: &Config{
					ServiceName: "test-service",
					ServicePort: "999999", // Invalid port
				},
			},
			before: nil, // Skip before since we expect a panic
			after:  nil, // No cleanup needed
			validate: func(t *testing.T, deps *serverTestDeps, args *serverTestArgs, err error) {
				// This test is expected to panic, so this should never be called
				t.Fatal("Test should have panicked before reaching validation")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "fails with invalid port" {
				// Special case for the panic test
				assert.Panics(t, func() {
					server := SetupServer(tt.args.config, &ApplicationComponents{
						OddsService:   &mockOddsService{},
						OddsValidator: validator.NewDefaultOddsValidator(NewLeagueRegistry(tt.args.config)),
					})
					server.Start()
				}, "Expected panic with invalid port")
				return
			}

			deps := &serverTestDeps{}
			args := tt.args

			// Run before hook if provided
			if tt.before != nil {
				tt.before(t, deps, &args)
			}

			// Ensure after hook runs
			if tt.after != nil {
				defer tt.after(t, deps, &args)
			}

			// Run validation if provided
			if tt.validate != nil {
				var err error
				if deps.server != nil && deps.server.grpcServer != nil {
					// If server started successfully, get any error from the error channel
					select {
					case err = <-deps.errChan:
					default:
						err = nil
					}
				}
				tt.validate(t, deps, &args, err)
			}
		})
	}
}

func TestRegisterWithGateway(t *testing.T) {
	// Create a log capture
	logCap := &logCapture{}
	log.SetOutput(logCap)
	defer log.SetOutput(io.Discard)

	// Start a mock gateway server on a random port
	mockGateway := &mockGatewayService{}
	gatewayServer := grpc.NewServer()
	proto.RegisterOddsServiceServer(gatewayServer, mockGateway)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer lis.Close()

	gatewayAddr := lis.Addr().String()

	// Start the mock gateway server in a goroutine
	go func() {
		if err := gatewayServer.Serve(lis); err != nil {
			t.Logf("Mock gateway server failed: %v", err)
		}
	}()
	defer gatewayServer.Stop()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	tests := []registerWithGatewayTest{
		{
			name: "successful registration",
			deps: registerWithGatewayTestDeps{
				gatewayServer: gatewayServer,
				gatewayAddr:   gatewayAddr,
			},
			args: registerWithGatewayTestArgs{
				gatewayAddr: gatewayAddr,
				servicePort: "50051",
			},
			before: func(t testing.TB, deps *registerWithGatewayTestDeps, args *registerWithGatewayTestArgs) {
				// Start the mock gateway server
				go func() {
					if err := deps.gatewayServer.Serve(lis); err != nil {
						t.Logf("Failed to serve: %v", err)
					}
				}()
			},
			after: func(t testing.TB, deps *registerWithGatewayTestDeps, args *registerWithGatewayTestArgs) {
				deps.gatewayServer.Stop()
			},
			expectedLogs: []string{
				"Attempting to register with API gateway at",
				"Successfully connected to API Gateway at",
			},
		},
		{
			name: "skip registration for default gateway address",
			deps: registerWithGatewayTestDeps{
				gatewayServer: gatewayServer,
			},
			args: registerWithGatewayTestArgs{
				gatewayAddr: "localhost:8080",
				servicePort: "50051",
			},
			after: func(t testing.TB, deps *registerWithGatewayTestDeps, args *registerWithGatewayTestArgs) {
				deps.gatewayServer.Stop()
			},
			expectedLogs: []string{
				"Default gateway address detected. Skipping automatic registration.",
				"Service is available at port 50051",
			},
		},
		{
			name: "connection failure",
			deps: registerWithGatewayTestDeps{
				gatewayServer: gatewayServer,
			},
			args: registerWithGatewayTestArgs{
				gatewayAddr: "localhost:9999", // Unavailable port
				servicePort: "50051",
			},
			after: func(t testing.TB, deps *registerWithGatewayTestDeps, args *registerWithGatewayTestArgs) {
				deps.gatewayServer.Stop()
			},
			expectedLogs: []string{
				"Attempting to register with API gateway at localhost:9999 (attempt 1/3)",
				"API Gateway connection test failed: rpc error: code = Unavailable desc = connection error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset log capture for this test case
			logCap.Reset()

			// Create test server with config
			cfg := &Config{
				ServiceName: "test-service",
				ServicePort: tt.args.servicePort,
				GatewayAddr: tt.args.gatewayAddr,
			}

			server := &Server{
				config: cfg,
			}

			// Run before hook if provided
			if tt.before != nil {
				tt.before(t, &tt.deps, &tt.args)
			}

			// Ensure after hook runs
			if tt.after != nil {
				defer tt.after(t, &tt.deps, &tt.args)
			}

			// Call the function under test
			server.RegisterWithGateway()

			// Wait for async operations to complete
			maxRetries := 10
			for range maxRetries {
				// Check if we've seen all expected logs
				allFound := true
				for _, expectedLog := range tt.expectedLogs {
					if !logCap.Contains(expectedLog) {
						allFound = false
						break
					}
				}
				if allFound {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			// Verify logs
			logs := logCap.String()
			for _, expectedLog := range tt.expectedLogs {
				assert.True(t, strings.Contains(logs, expectedLog),
					"Expected log message containing: %s\nGot logs: %s", expectedLog, logs)
			}

			// Run validation if provided
			if tt.validate != nil {
				tt.validate(t, &tt.deps, &tt.args, server)
			}
		})
	}
}

// logCapture is a simple writer that captures log output for testing
type logCapture struct {
	mu   sync.Mutex
	logs []byte
}

// Write implements io.Writer interface
func (l *logCapture) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = append(l.logs, p...)
	return len(p), nil
}

// String returns the captured log output as a string
func (l *logCapture) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return string(l.logs)
}

// Contains checks if the logs contain the given string
func (l *logCapture) Contains(s string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Contains(string(l.logs), s)
}

// Reset clears the captured logs
func (l *logCapture) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = nil
}

type mockOddsService struct {
	ports.OddsService
}

// Ensure mockOddsService implements ports.OddsService
var _ ports.OddsService = (*mockOddsService)(nil)

func testHealthCheck(t *testing.T, addr string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	resp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	require.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
}

type mockGatewayService struct {
	proto.UnimplementedOddsServiceServer
	shouldFail bool
}

func (m *mockGatewayService) GetOdds(ctx context.Context, req *proto.GetOddsRequest) (*proto.GetOddsResponse, error) {
	if m.shouldFail {
		return nil, assert.AnError
	}
	return &proto.GetOddsResponse{}, nil
}

type registerWithGatewayTestDeps struct {
	gatewayServer *grpc.Server
	gatewayAddr   string
}

type registerWithGatewayTestArgs struct {
	gatewayAddr string
	servicePort string
}

type registerWithGatewayTest struct {
	name         string
	deps         registerWithGatewayTestDeps
	args         registerWithGatewayTestArgs
	before       func(testing.TB, *registerWithGatewayTestDeps, *registerWithGatewayTestArgs)
	after        func(testing.TB, *registerWithGatewayTestDeps, *registerWithGatewayTestArgs)
	validate     func(testing.TB, *registerWithGatewayTestDeps, *registerWithGatewayTestArgs, *Server)
	expectedLogs []string
}
