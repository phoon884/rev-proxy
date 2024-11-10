package application

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/phoon884/rev-proxy/pkg/logger/ports"
)

type Logger struct {
	zap *zap.SugaredLogger
}

var _ ports.LoggerApplication = (*Logger)(nil)

func NewLogger(logLevel string) *Logger {
	cfg := zap.NewProductionConfig()
	var zapLevel zapcore.Level
	found := true
	switch logLevel {
	case "DEBUG":
		zapLevel = zap.DebugLevel
		break
	case "INFO":
		zapLevel = zap.InfoLevel
		break
	case "WARN":
		zapLevel = zap.WarnLevel
		break
	case "ERROR":
		zapLevel = zap.ErrorLevel
		break
	default:
		found = false
		zapLevel = zap.InfoLevel
	}

	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	logger, _ := cfg.Build()
	sugar := logger.Sugar()
	if !found {
		sugar.Warn("Log level is not inputted...")
	}
	sugar.Info("Log level = INFO")

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
