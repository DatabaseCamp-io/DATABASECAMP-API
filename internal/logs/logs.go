package logs

// logs.go
/**
 * 	This file used to manage logs of the application
 */

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/**
 * This class create manage of the application
 */
type log struct {
	log *zap.Logger
}

// Instance of log class for singleton pattern
var instantiated *log = nil

/**
 * Constructor creates a new log instance or geting a log instance
 *
 * @return 	instance of log
 */
func GetInstance() *log {
	var buildedLog *zap.Logger
	var err error

	if instantiated == nil {
		buildedLog, err = initLog()
		if err != nil {
			panic(err)
		}
		instantiated = &log{log: buildedLog}
	}

	return instantiated
}

/**
 * Config log
 *
 * @return logger of the application
 * @return the error of building logger
 */
func initLog() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	buildedLog, err := config.Build(zap.AddCallerSkip(1))
	return buildedLog, err
}

/**
 * Log information
 *
 * @param message infomation log message
 * @param fields  optional for used to add a key-value pair to a logger's context
 */
func (l log) Info(message interface{}) {
	l.log.Info(fmt.Sprint(message))
}

/**
 * Log debug
 *
 * @param message	debug log message
 * @param fields  	optional for used to add a key-value pair to a logger's context
 */
func (l log) Debug(message string) {
	l.log.Debug(message)
}

/**
 * Log error
 *
 * @param message	error log message
 * @param fields  	optional for used to add a key-value pair to a logger's context
 */
func (l log) Error(message interface{}) {
	switch v := message.(type) {
	case error:
		l.log.Error(v.Error())
	case string:
		l.log.Error(v)
	}
}
