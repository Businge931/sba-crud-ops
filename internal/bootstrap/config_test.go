package bootstrap

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDummy(t *testing.T) {
	assert.True(t, true, "This test should always pass")
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
		expectedErr    bool
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expectedConfig: &Config{
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
			},
		},
		{
			name: "custom values",
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
			// Setup
			resetEnvVars()
			setEnvVars(t, tt.envVars)

			// Execute
			config := LoadConfig()

			// Verify
			assertConfigEqual(t, tt.expectedConfig, config)
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
