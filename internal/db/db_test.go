package db_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/bootstrap"
	"github.com/Businge931/sba-crud-ops/internal/db"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSetupGormDB(t *testing.T) {
	type dependencies struct {
		container testcontainers.Container
	}
	type args struct {
		cfg *bootstrap.Config
	}
	testCases := []struct {
		name         string
		dependencies dependencies
		args         args
		before       func(t *testing.T, deps *dependencies, args *args)
		after        func(t *testing.T, deps *dependencies, args *args)
		want         any
		wantErr      bool
	}{
		{
			name:         "valid config (containerized test DB)",
			dependencies: dependencies{},
			args:         args{cfg: nil},
			before: func(t *testing.T, deps *dependencies, args *args) {
				// Start testcontainer for PostgreSQL
				testUser := "testuser"
				testPass := "testpass"
				testDB := "testdb"

				ctx := context.Background()
				containerReq := testcontainers.ContainerRequest{
					Image:        "postgres:15-alpine",
					ExposedPorts: []string{"5432/tcp"},
					Env: map[string]string{
						"POSTGRES_USER":     testUser,
						"POSTGRES_PASSWORD": testPass,
						"POSTGRES_DB":       testDB,
					},
					WaitingFor: wait.ForLog("database system is ready to accept connections"),
				}
				postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
					ContainerRequest: containerReq,
					Started:          true,
				})
				if err != nil {
					t.Skip("Docker/testcontainers not available: ", err)
				}
				// Get host/port
				host, err := postgresC.Host(ctx)
				if err != nil {
					postgresC.Terminate(ctx)
					t.Fatalf("failed to get container host: %v", err)
				}
				port, err := postgresC.MappedPort(ctx, "5432")
				if err != nil {
					postgresC.Terminate(ctx)
					t.Fatalf("failed to get container port: %v", err)
				}
				args.cfg = &bootstrap.Config{
					DBAddr:       "host=" + host + " port=" + port.Port() + " user=" + testUser + " password=" + testPass + " dbname=" + testDB + " sslmode=disable",
					MaxIdleTime:  "5m",
					MaxOpenConns: 10,
					MaxIdleConns: 5,
				}
				// Wait for DB port to be ready
				ready := false
				addr := host + ":" + port.Port()
				for i := 0; i < 50; i++ { // 50 * 200ms = 10s
					conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
					if err == nil {
						conn.Close()
						ready = true
						break
					}
					time.Sleep(200 * time.Millisecond)
				}
				if !ready {
					postgresC.Terminate(ctx)
					t.Fatalf("database port %s not ready after 10s", addr)
				}
				// Save container for cleanup
				deps.container = postgresC
			},
			after: func(t *testing.T, deps *dependencies, args *args) {
				if deps.container != nil {
					ctx := context.Background()
					_ = deps.container.Terminate(ctx)
				}
			},
			want:    "*gorm.DB",
			wantErr: false,
		},
		{
			name: "invalid connection string",
			args: args{
				cfg: &bootstrap.Config{
					DBAddr: "invalid-connection-string",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty connection string",
			args: args{
				cfg: &bootstrap.Config{
					DBAddr: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}
			var (
				dbConn *gorm.DB
				dbErr  error
			)
			if tc.name == "valid config (containerized test DB)" {
				t.Logf("Connecting with DSN: %s", tc.args.cfg.DBAddr)
				for i := 0; i < 10; i++ {
					dbConn, dbErr = db.SetupGormDB(tc.args.cfg)
					if dbErr == nil {
						break
					}
					t.Logf("Retrying SetupGormDB after error: %v", dbErr)
					time.Sleep(1 * time.Second)
				}
				if tc.wantErr {
					require.Error(t, dbErr)
					require.Nil(t, dbConn)
				} else {
					require.NoError(t, dbErr)
					require.IsType(t, &gorm.DB{}, dbConn)
					sqlDB, err := dbConn.DB()
					require.NoError(t, err)
					require.NoError(t, sqlDB.PingContext(context.Background()))
				}
			} else {
				db, err := db.SetupGormDB(tc.args.cfg)
				if tc.wantErr {
					require.Error(t, err)
					require.Nil(t, db)
				} else {
					require.NoError(t, err)
					require.IsType(t, &gorm.DB{}, db)
					sqlDB, err := db.DB()
					require.NoError(t, err)
					require.NoError(t, sqlDB.PingContext(context.Background()))
				}
			}
			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
		})
	}
}

// TestRunMigrations is a standardized, table-driven integration test for db.RunMigrations.
// It uses testcontainers to spin up a real Postgres DB and tests both success and error cases.
func TestRunMigrations(t *testing.T) {
	type dependencies struct {
		container testcontainers.Container
	}
	type args struct {
		db *gorm.DB
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
			name:         "success (containerized test DB)",
			dependencies: dependencies{},
			args:         args{db: nil},
			before: func(t *testing.T, deps *dependencies, args *args) {
				// Start testcontainer for PostgreSQL
				testUser := "testuser"
				testPass := "testpass"
				testDB := "testdb"

				ctx := context.Background()
				containerReq := testcontainers.ContainerRequest{
					Image:        "postgres:15-alpine",
					ExposedPorts: []string{"5432/tcp"},
					Env: map[string]string{
						"POSTGRES_USER":     testUser,
						"POSTGRES_PASSWORD": testPass,
						"POSTGRES_DB":       testDB,
					},
					WaitingFor: wait.ForLog("database system is ready to accept connections"),
				}
				postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
					ContainerRequest: containerReq,
					Started:          true,
				})
				if err != nil {
					t.Skip("Docker/testcontainers not available: ", err)
				}
				// Get host/port
				host, err := postgresC.Host(ctx)
				if err != nil {
					postgresC.Terminate(ctx)
					t.Fatalf("failed to get container host: %v", err)
				}
				port, err := postgresC.MappedPort(ctx, "5432")
				if err != nil {
					postgresC.Terminate(ctx)
					t.Fatalf("failed to get container port: %v", err)
				}
				cfg := &bootstrap.Config{
					DBAddr:       "host=" + host + " port=" + port.Port() + " user=" + testUser + " password=" + testPass + " dbname=" + testDB + " sslmode=disable",
					MaxIdleTime:  "5m",
					MaxOpenConns: 10,
					MaxIdleConns: 5,
				}
				// Wait for DB port to be ready
				ready := false
				addr := host + ":" + port.Port()
				for i := 0; i < 50; i++ {
					conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
					if err == nil {
						conn.Close()
						ready = true
						break
					}
					time.Sleep(200 * time.Millisecond)
				}
				if !ready {
					postgresC.Terminate(ctx)
					t.Fatalf("database port %s not ready after 10s", addr)
				}
				// Setup GORM DB (with retry)
				var dbConn *gorm.DB
				var dbErr error
				for i := 0; i < 10; i++ {
					dbConn, dbErr = db.SetupGormDB(cfg)
					if dbErr == nil {
						break
					}
					t.Logf("Retrying SetupGormDB after error: %v", dbErr)
					time.Sleep(1 * time.Second)
				}
				if dbErr != nil {
					postgresC.Terminate(ctx)
					t.Fatalf("failed to setup GORM DB: %v", dbErr)
				}
				args.db = dbConn
				deps.container = postgresC
			},
			after: func(t *testing.T, deps *dependencies, args *args) {
				if deps.container != nil {
					ctx := context.Background()
					_ = deps.container.Terminate(ctx)
				}
			},
			wantErr: false,
		},
		{
			name:    "nil db returns error",
			args:    args{db: nil},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}
			err := db.RunMigrations(tc.args.db)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
		})
	}
}
