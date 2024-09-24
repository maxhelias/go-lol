package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogField is an alias of zapcore.Field, it could be replaced by a custom contract when Go will support generics.
type LogField = zapcore.Field

// Level is an alias of zapcore.Level, it could be replaced by a custom contract when Go will support generics.
type Level = zapcore.Level

// CheckedEntry is an alias of zapcore.CheckedEntry, it could be replaced by a custom contract when Go will support generics.
type CheckedEntry = zapcore.CheckedEntry

type Logger interface {
	Info(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
	Check(level Level, msg string) *CheckedEntry
	Level() Level
}

type ZapLogger struct {
	logger *zap.Logger
}

// NewLogger creates a new logger.
func NewLogger(debug bool) (*ZapLogger, error) {
	var zapLogger *zap.Logger
	var err error
	if debug {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: zapLogger}, nil
}

func (zl *ZapLogger) Info(msg string, fields ...zapcore.Field) {
	zl.logger.Info(msg, fields...)
}

func (zl *ZapLogger) Error(msg string, fields ...zapcore.Field) {
	zl.logger.Error(msg, fields...)
}

func (zl *ZapLogger) Check(level zapcore.Level, msg string) *zapcore.CheckedEntry {
	return zl.logger.Check(level, msg)
}

func (zl *ZapLogger) Level() zapcore.Level {
	return zl.logger.Level()
}
