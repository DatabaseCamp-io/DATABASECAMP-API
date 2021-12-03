package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type log struct {
	log *zap.Logger
}

var instantiated *log = nil

// Create log instance
func New() *log {
	var buildedLog *zap.Logger
	if instantiated == nil {
		buildedLog, _ = initLog()
		instantiated = &log{log: buildedLog}
	}
	return instantiated
}

// Init log
func initLog() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	buildedLog, err := config.Build(zap.AddCallerSkip(1))
	return buildedLog, err
}

// Create info logs
func (l log) Info(message string, fields ...zapcore.Field) {
	l.log.Info(message, fields...)
}

// Debug logs
func (l log) Debug(message string, fields ...zapcore.Field) {
	l.log.Debug(message, fields...)
}

// Return error message log
func (l log) Error(message interface{}, fields ...zapcore.Field) {
	switch v := message.(type) {
	case error:
		l.log.Error(v.Error(), fields...)
	case string:
		l.log.Error(v, fields...)
	}
}
