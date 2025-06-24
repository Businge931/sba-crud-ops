package usecase

import (
	"maps"
	"time"

	log "github.com/sirupsen/logrus"
)

// BaseLogger provides common logging functionality for all use cases
type BaseLogger struct{}

// LogStart logs the beginning of an operation with the given parameters
func (bl *BaseLogger) LogStart(operation string, fields log.Fields) {
	fields["operation"] = operation
	log.WithFields(fields).Info("Starting operation")
}

// LogSuccess logs the successful completion of an operation
func (bl *BaseLogger) LogSuccess(operation string, startTime time.Time, extraFields log.Fields) {
	fields := log.Fields{
		"operation": operation,
		"duration":  time.Since(startTime).String(),
	}

	// Add any extra fields
	maps.Copy(fields, extraFields)

	log.WithFields(fields).Info("Successfully completed operation")
}

// LogError logs a failed operation with error details
func (bl *BaseLogger) LogError(operation string, startTime time.Time, err error) {
	log.WithFields(log.Fields{
		"error":     err.Error(),
		"operation": operation,
		"duration":  time.Since(startTime).String(),
	}).Error("Operation failed")
}
