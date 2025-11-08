package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is a wrapper around zap.Logger that implements the Logger interface
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new ZapLogger with the given mode
func NewZapLogger(mode string) Logger {
	var zapConfig zap.Config

	if mode == "release" {
		// Production mode: JSON encoding, Info level
		// Optimized for Docker/Cloud - logs to stdout/stderr only
		zapConfig = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		zapConfig.EncoderConfig.StacktraceKey = ""
	} else {
		// Development mode: Console encoding, Debug level
		zapConfig = zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		// Customize encoder config for better readability
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	logger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		logger, _ = zap.NewProduction()
	}

	return &ZapLogger{logger: logger}
}

// fieldsToZap converts our generic Field to zap.Field
func fieldsToZap(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fieldsToZap(fields)...)
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fieldsToZap(fields)...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fieldsToZap(fields)...)
}

func (l *ZapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fieldsToZap(fields)...)
}

func (l *ZapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fieldsToZap(fields)...)
	os.Exit(1)
}

// Sync flushes any buffered log entries
func (l *ZapLogger) Sync() {
	_ = l.logger.Sync()
}

// GetZapLogger returns the underlying zap.Logger for advanced usage
// This is useful for middleware that requires *zap.Logger directly
func (l *ZapLogger) GetZapLogger() *zap.Logger {
	return l.logger
}
