package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func InitLogger(level, format string) {
	var lvl zapcore.Level
	switch level {
	case "debug":
		lvl = zapcore.DebugLevel
	case "warn":
		lvl = zapcore.WarnLevel
	case "error":
		lvl = zapcore.ErrorLevel
	default:
		lvl = zapcore.InfoLevel
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.Encoding = "json"
	if format == "console" {
		cfg.Encoding = "console"
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic("Failed to init zap logger: " + err.Error())
	}
	zap.ReplaceGlobals(log)
}

func L() *zap.Logger {
	if log == nil {
		return zap.L()
	}
	return log
}

func S() *zap.SugaredLogger {
	return L().Sugar()
}
