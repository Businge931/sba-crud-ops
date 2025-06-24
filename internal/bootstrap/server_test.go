package bootstrap

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/adoptors/secondary/validator"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
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
	type dependencies struct{}
	type args struct {
		cfg        *Config
		components *ApplicationComponents
	}

	testCases := []struct {
		name         string
		dependencies dependencies
		args         args
		before       func(t *testing.T, deps *dependencies, args *args)
		after        func(t *testing.T, deps *dependencies, args *args)
		wantErr      bool
	}{
		{
			name:         "server setup and start",
			dependencies: dependencies{},
			args: args{
				cfg: &Config{
					ServicePort: "0",
					ServiceName: "test-service",
				},
				components: &ApplicationComponents{
					OddsService:   &mockOddsService{},
					OddsValidator: validator.NewDefaultOddsValidator(NewLeagueRegistry(&Config{ServiceName: "test-service"})),
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}
			server := SetupServer(tc.args.cfg, tc.args.components)
			require.NotNil(t, server)
			require.NotNil(t, server.grpcServer)

			errChan := make(chan error, 1)
			go func() {
				errChan <- server.Start()
			}()
			time.Sleep(100 * time.Millisecond)
			testHealthCheck(t, server.listener.Addr().String())
			server.grpcServer.GracefulStop()
			select {
			case err := <-errChan:
				assert.NoError(t, err)
			case <-time.After(1 * time.Second):
				t.Fatal("Server did not stop within timeout")
			}
			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
		})
	}
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
	type dependencies struct{}
	type args struct {
		config *Config
	}

	testCases := []struct {
		name         string
		dependencies dependencies
		args         args
		before       func(t *testing.T, deps *dependencies, args *args)
		after        func(t *testing.T, deps *dependencies, args *args)
		wantErr      bool
	}{
		{
			name:         "successful registration",
			dependencies: dependencies{},
			args:         args{config: &Config{ServiceName: "test-service", GatewayAddr: "localhost:8080"}},
			wantErr:      false,
		},
		{
			name:         "invalid gateway address",
			dependencies: dependencies{},
			args:         args{config: &Config{ServiceName: "test-service", GatewayAddr: "invalid-address"}},
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}
			server := &Server{
				config:     tc.args.config,
				grpcServer: grpc.NewServer(),
			}
			lis, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			server.listener = lis
			t.Skip("Skipping RegisterWithGateway test as it requires a running gateway server")
			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
		})
	}
}

func TestServerStartError(t *testing.T) {
	type dependencies struct{}
	type args struct {
		config *Config
	}

	testCases := []struct {
		name         string
		dependencies dependencies
		args         args
		before       func(t *testing.T, deps *dependencies, args *args)
		after        func(t *testing.T, deps *dependencies, args *args)
		wantPanic    bool
	}{
		{
			name:         "invalid port triggers panic",
			dependencies: dependencies{},
			args:         args{config: &Config{ServicePort: "99999"}},
			wantPanic:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}
			server := &Server{
				config:     tc.args.config,
				grpcServer: grpc.NewServer(),
			}
			if tc.wantPanic {
				assert.Panics(t, func() {
					err := server.Start()
					if err != nil {
						panic(err)
					}
				}, "Expected Start to panic with invalid port")
			} else {
				assert.NotPanics(t, func() {
					_ = server.Start()
				})
			}
			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
		})
	}
}
