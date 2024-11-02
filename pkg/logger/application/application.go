package application

import (
	"go.uber.org/zap"

	"github.com/phoon884/rev-proxy/pkg/logger/ports"
)

type Logger struct {
	zap *zap.SugaredLogger
}

var _ ports.LoggerApplication = (*Logger)(nil)

func NewLogger() *Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, _ := cfg.Build()
	sugar := logger.Sugar()

	return &Logger{zap: sugar}
}

func (l Logger) Close() error {
	return l.zap.Sync()
}

func (l Logger) Debug(msg string, args ...interface{}) {
	l.zap.Debug(msg, args)
}

func (l Logger) Info(msg string, args ...interface{}) {
	l.zap.Info(msg, args)
}

func (l Logger) Warn(msg string, args ...interface{}) {
	l.zap.Warn(msg, args)
}

func (l Logger) Error(msg string, args ...interface{}) {
	l.zap.Error(msg, args)
}
