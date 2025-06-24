package usecase

import (
	"maps"
	"time"

	"github.com/sirupsen/logrus"
)

// BaseLogger provides common logging functionality for use cases
type BaseLogger struct {
	logger *logrus.Logger
}

// NewBaseLogger creates a new BaseLogger with the given logrus.Logger
func NewBaseLogger(logger *logrus.Logger) BaseLogger {
	if logger == nil {
		return BaseLogger{logger: logrus.StandardLogger()}
	}
	return BaseLogger{logger: logger}
}

// LogStart logs the beginning of an operation with the given parameters
func (bl *BaseLogger) LogStart(operation string, fields logrus.Fields) {
	fields["operation"] = operation
	if bl.logger != nil {
		bl.logger.WithFields(fields).Info("Starting operation")
	} else {
		logrus.WithFields(fields).Info("Starting operation")
	}
}

// LogSuccess logs the successful completion of an operation
func (bl *BaseLogger) LogSuccess(operation string, startTime time.Time, extraFields logrus.Fields) {
	fields := logrus.Fields{
		"operation": operation,
		"duration":  time.Since(startTime).String(),
	}

	// Add any extra fields
	if extraFields != nil {
		maps.Copy(fields, extraFields)
	}

	if bl.logger != nil {
		bl.logger.WithFields(fields).Info("Successfully completed operation")
	} else {
		logrus.WithFields(fields).Info("Successfully completed operation")
	}
}

// LogError logs a failed operation with error details
func (bl *BaseLogger) LogError(operation string, startTime time.Time, err error) {
	fields := logrus.Fields{
		"error":     err.Error(),
		"operation": operation,
		"duration":  time.Since(startTime).String(),
	}
	if bl.logger != nil {
		bl.logger.WithFields(fields).Error("Operation failed")
	} else {
		logrus.WithFields(fields).Error("Operation failed")
	}
}
