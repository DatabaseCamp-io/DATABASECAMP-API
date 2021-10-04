package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	log *zap.Logger
}

type ILog interface {
	Info(message string, fields ...zapcore.Field)
	Debug(message string, fields ...zapcore.Field)
	Error(message interface{}, fields ...zapcore.Field)
}

var instantiated *Log = nil

func New() *Log {
	if instantiated == nil {
		instantiated = new(Log)
		instantiated.init()
	}
	return instantiated
}

func (l Log) init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	buildedLog, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	l.log = buildedLog
}

func (l Log) Info(message string, fields ...zapcore.Field) {
	l.log.Info(message, fields...)
}

func (l Log) Debug(message string, fields ...zapcore.Field) {
	l.log.Debug(message, fields...)
}

func (l Log) Error(message interface{}, fields ...zapcore.Field) {
	switch v := message.(type) {
	case error:
		l.log.Error(v.Error(), fields...)
	case string:
		l.log.Error(v, fields...)
	}
}
