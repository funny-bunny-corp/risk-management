package logger

import "go.uber.org/zap"

func NewLogger() *zap.Logger {
	logger := zap.Must(zap.NewProduction())
	return logger
}
