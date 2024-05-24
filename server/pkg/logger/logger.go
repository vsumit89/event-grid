package logger

import (
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

type logFmt string

const (
	JSON logFmt = "JSON"
	TEXT logFmt = "TEXT"
)

type loggerOptionFunc func(*loggerConfig)

type LOG_LEVEL string

const (
	INFO  LOG_LEVEL = "INFO"
	WARN  LOG_LEVEL = "WARN"
	DEBUG LOG_LEVEL = "DEBUG"
	ERROR LOG_LEVEL = "ERROR"
)

func getLogLevel(level LOG_LEVEL) slog.Level {
	switch level {
	case INFO:
		return slog.LevelInfo
	case WARN:
		return slog.LevelWarn
	case DEBUG:
		return slog.LevelDebug
	case ERROR:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type loggerConfig struct {
	level  LOG_LEVEL
	writer io.Writer
	logFmt logFmt
}

type LoggerOption interface {
	apply(*loggerConfig)
}

func (f loggerOptionFunc) apply(config *loggerConfig) {
	f(config)
}

func WithLevel(level LOG_LEVEL) LoggerOption {
	return loggerOptionFunc(func(config *loggerConfig) {
		config.level = level
	})
}

func WithWriter(writer io.Writer) LoggerOption {
	return loggerOptionFunc(func(config *loggerConfig) {
		config.writer = writer
	})
}

func WithLogFmt(logFmt logFmt) LoggerOption {
	return loggerOptionFunc(func(config *loggerConfig) {
		config.logFmt = logFmt
	})
}

var logger = &Logger{}

func InitLogger(options ...LoggerOption) {
	defaultCfg := &loggerConfig{
		level:  DEBUG,
		writer: os.Stdout,
		logFmt: TEXT,
	}

	for _, opt := range options {
		// applying options passed by the caller
		opt.apply(defaultCfg)
	}

	var loggerInstance *slog.Logger

	if defaultCfg.logFmt == JSON {
		loggerInstance = slog.New(slog.NewJSONHandler(defaultCfg.writer, &slog.HandlerOptions{
			Level: getLogLevel(defaultCfg.level),
		}))
	} else {
		loggerInstance = slog.New(slog.NewTextHandler(defaultCfg.writer, &slog.HandlerOptions{
			Level: getLogLevel(defaultCfg.level),
		}))
	}

	logger.logger = loggerInstance

}

func Info(message string, v ...interface{}) {
	logger.logger.Info(message, v...)
}

func Warn(message string, v ...interface{}) {
	logger.logger.Warn(message, v...)
}

func Debug(message string, v ...interface{}) {
	logger.logger.Debug(message, v...)
}

func Error(message string, v ...interface{}) {
	logger.logger.Error(message, v...)
}
