package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

type baseLoggerTestDeps struct {
	hook *test.Hook
}

type baseLoggerTestArgs struct {
	operation   string
	fields      logrus.Fields
	startTime   time.Time
	extraFields logrus.Fields
	err         error
}

type baseLoggerTestExpected struct {
	message    string
	level      logrus.Level
	dataFields map[string]interface{}
}

func TestBaseLogger_LogStart(t *testing.T) {
	tests := []struct {
		name     string
		deps     func() *baseLoggerTestDeps
		args     baseLoggerTestArgs
		expected baseLoggerTestExpected
	}{
		{
			name: "should log operation start with fields",
			deps: func() *baseLoggerTestDeps {
				hook := test.NewGlobal()
				logrus.SetLevel(logrus.InfoLevel)
				logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
				return &baseLoggerTestDeps{hook: hook}
			},
			args: baseLoggerTestArgs{
				operation: "testOperation",
				fields:    logrus.Fields{"key1": "value1", "key2": 42},
			},
			expected: baseLoggerTestExpected{
				message: "Starting operation",
				level:   logrus.InfoLevel,
				dataFields: map[string]interface{}{
					"operation": "testOperation",
					"key1":      "value1",
					"key2":      42,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			deps := tt.deps()
			logger := &BaseLogger{}

			// Execute
			logger.LogStart(tt.args.operation, tt.args.fields)

			// Verify
			assert.Equal(t, 1, len(deps.hook.Entries))
			entry := deps.hook.LastEntry()
			assert.Equal(t, tt.expected.message, entry.Message)
			assert.Equal(t, tt.expected.level, entry.Level)

			for k, v := range tt.expected.dataFields {
				assert.Equal(t, v, entry.Data[k])
			}

			// Cleanup
			deps.hook.Reset()
		})
	}
}

func TestBaseLogger_LogSuccess(t *testing.T) {
	tests := []struct {
		name     string
		deps     func() *baseLoggerTestDeps
		args     baseLoggerTestArgs
		expected baseLoggerTestExpected
	}{
		{
			name: "should log operation success with duration",
			deps: func() *baseLoggerTestDeps {
				hook := test.NewGlobal()
				logrus.SetLevel(logrus.InfoLevel)
				return &baseLoggerTestDeps{hook: hook}
			},
			args: baseLoggerTestArgs{
				operation:   "testOperation",
				extraFields: logrus.Fields{"result": "success", "count": 1},
				startTime:   time.Now().Add(-100 * time.Millisecond),
			},
			expected: baseLoggerTestExpected{
				message: "Successfully completed operation",
				level:   logrus.InfoLevel,
				dataFields: map[string]interface{}{
					"operation": "testOperation",
					"result":    "success",
					"count":     1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			deps := tt.deps()
			logger := &BaseLogger{}

			// Execute
			logger.LogSuccess(tt.args.operation, tt.args.startTime, tt.args.extraFields)

			// Verify
			assert.Equal(t, 1, len(deps.hook.Entries))
			entry := deps.hook.LastEntry()
			assert.Equal(t, tt.expected.message, entry.Message)
			assert.Equal(t, tt.expected.level, entry.Level)

			// Check duration is present and positive
			assert.Contains(t, entry.Data, "duration")
			durationStr, ok := entry.Data["duration"].(string)
			assert.True(t, ok)
			duration, err := time.ParseDuration(durationStr)
			assert.NoError(t, err)
			assert.True(t, duration > 0)

			// Check other fields
			for k, v := range tt.expected.dataFields {
				assert.Equal(t, v, entry.Data[k])
			}

			// Cleanup
			deps.hook.Reset()
		})
	}
}

func TestBaseLogger_LogError(t *testing.T) {
	tests := []struct {
		name     string
		deps     func() *baseLoggerTestDeps
		args     baseLoggerTestArgs
		expected baseLoggerTestExpected
	}{
		{
			name: "should log operation failure with error",
			deps: func() *baseLoggerTestDeps {
				hook := test.NewGlobal()
				logrus.SetLevel(logrus.ErrorLevel)
				return &baseLoggerTestDeps{hook: hook}
			},
			args: baseLoggerTestArgs{
				operation: "testOperation",
				err:       errors.New("test error"),
				startTime: time.Now().Add(-150 * time.Millisecond),
			},
			expected: baseLoggerTestExpected{
				message: "Operation failed",
				level:   logrus.ErrorLevel,
				dataFields: map[string]interface{}{
					"operation": "testOperation",
					"error":     "test error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			deps := tt.deps()
			logger := &BaseLogger{}

			// Execute
			logger.LogError(tt.args.operation, tt.args.startTime, tt.args.err)

			// Verify
			assert.Equal(t, 1, len(deps.hook.Entries))
			entry := deps.hook.LastEntry()
			assert.Equal(t, tt.expected.message, entry.Message)
			assert.Equal(t, tt.expected.level, entry.Level)

			// Check duration is present and positive
			assert.Contains(t, entry.Data, "duration")
			durationStr, ok := entry.Data["duration"].(string)
			assert.True(t, ok)
			duration, err := time.ParseDuration(durationStr)
			assert.NoError(t, err)
			assert.True(t, duration > 0)

			// Check other fields
			for k, v := range tt.expected.dataFields {
				assert.Equal(t, v, entry.Data[k])
			}

			// Cleanup
			deps.hook.Reset()
		})
	}
}
