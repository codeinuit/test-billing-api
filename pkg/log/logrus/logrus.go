package logrus

import "github.com/sirupsen/logrus"

type LogrusLogger struct {
	logger *logrus.Logger
}

// NewLogrusLogger returns a log instance from
// logger implementation
func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()
	return &LogrusLogger{
		logger: logger,
	}
}

// Info prints logs at info level
func (l LogrusLogger) Info(v ...any) {
	l.logger.Info(v...)
}

// Error prints logs at error level
func (l LogrusLogger) Error(v ...any) {
	l.logger.Error(v...)
}

// Warn prints logs at warn level
func (l LogrusLogger) Warn(v ...any) {
	l.logger.Warn(v...)
}
