package bootstrap

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDummy(t *testing.T) {
	assert.True(t, true, "This test should always pass")
}

type loadConfigTestDeps struct {
	// Currently no external dependencies for this test
}

type loadConfigTestArgs struct {
	envVars map[string]string
}

type loadConfigTest struct {
	name           string
	deps           loadConfigTestDeps
	args           loadConfigTestArgs
	after          func(testing.TB, *loadConfigTestDeps)
	expectedConfig *Config
}

func TestLoadConfig(t *testing.T) {
	defaultConfig := &Config{
		ServicePort:  "50052",
		ServiceName:  "odds-service",
		Version:      "0.0.1",
		GatewayAddr:  "localhost:8080",
		DBAddr:       "postgresql://admin:adminpassword@localhost:5433/sba_crud_ops?sslmode=disable",
		MaxOpenConns: 30,
		MaxIdleConns: 30,
		MaxIdleTime:  "15m",
		SupportedLeagues: []string{
			"English Premier League",
			"La Liga",
			"Serie A",
			"Bundesliga",
			"Ligue 1",
		},
	}

	tests := []loadConfigTest{
		{
			name: "default values",
			args: loadConfigTestArgs{
				envVars: map[string]string{},
			},
			after: func(t testing.TB, _ *loadConfigTestDeps) {
				// Clean up environment variables after test
				resetEnvVars()
			},
			expectedConfig: defaultConfig,
		},
		{
			name: "custom values",
			args: loadConfigTestArgs{
				envVars: map[string]string{
					"SERVICE_PORT":      "8080",
					"SERVICE_NAME":      "test-service",
					"VERSION":           "1.0.0",
					"GATEWAY_ADDR":      "gateway:9090",
					"DB_ADDR":           "postgresql://user:pass@localhost:5432/testdb",
					"DB_MAX_OPEN_CONNS": "50",
					"DB_MAX_IDLE_CONNS": "25",
					"DB_MAX_IDLE_TIME":  "10m",
				},
			},
			after: func(t testing.TB, _ *loadConfigTestDeps) {
				// Clean up environment variables after test
				resetEnvVars()
			},
			expectedConfig: &Config{
				ServicePort:  "8080",
				ServiceName:  "test-service",
				Version:      "1.0.0",
				GatewayAddr:  "gateway:9090",
				DBAddr:       "postgresql://user:pass@localhost:5432/testdb",
				MaxOpenConns: 50,
				MaxIdleConns: 25,
				MaxIdleTime:  "10m",
				SupportedLeagues: []string{
					"English Premier League",
					"La Liga",
					"Serie A",
					"Bundesliga",
					"Ligue 1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			resetEnvVars()
			setEnvVars(t, tt.args.envVars)

			// Execute test
			config := LoadConfig()

			// Verify results
			assertConfigEqual(t, tt.expectedConfig, config)

			// Run after function if provided
			if tt.after != nil {
				tt.after(t, &tt.deps)
			}
		})
	}
}

// assertConfigEqual is a helper to assert config fields
func assertConfigEqual(t *testing.T, expected, actual *Config) {
	t.Helper()
	assert.Equal(t, expected.ServicePort, actual.ServicePort)
	assert.Equal(t, expected.ServiceName, actual.ServiceName)
	assert.Equal(t, expected.Version, actual.Version)
	assert.Equal(t, expected.GatewayAddr, actual.GatewayAddr)
	assert.Equal(t, expected.DBAddr, actual.DBAddr)
	assert.Equal(t, expected.MaxOpenConns, actual.MaxOpenConns)
	assert.Equal(t, expected.MaxIdleConns, actual.MaxIdleConns)
	assert.Equal(t, expected.MaxIdleTime, actual.MaxIdleTime)
	assert.ElementsMatch(t, expected.SupportedLeagues, actual.SupportedLeagues)
}

func resetEnvVars() {
	envVars := []string{
		"SERVICE_PORT",
		"SERVICE_NAME",
		"VERSION",
		"GATEWAY_ADDR",
		"DB_ADDR",
		"DB_MAX_OPEN_CONNS",
		"DB_MAX_IDLE_CONNS",
		"DB_MAX_IDLE_TIME",
	}

	for _, v := range envVars {
		os.Unsetenv(v)
	}
}

func setEnvVars(t *testing.T, envs map[string]string) {
	t.Helper()
	for k, v := range envs {
		t.Setenv(k, v)
	}
}
