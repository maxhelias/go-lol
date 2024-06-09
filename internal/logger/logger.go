package logger

import (
	"go.uber.org/zap/zapcore"
)

// LogField is an alias of zapcore.Field, it could be replaced by a custom contract when Go will support generics.
type LogField = zapcore.Field

// Level is an alias of zapcore.Level, it could be replaced by a custom contract when Go will support generics.
type Level = zapcore.Level

// CheckedEntry is an alias of zapcore.CheckedEntry, it could be replaced by a custom contract when Go will support generics.
type CheckedEntry = zapcore.CheckedEntry

// Logger defines the Mercure logger.
type Logger interface {
	Info(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
	Check(level Level, msg string) *CheckedEntry
	Level() Level
}
