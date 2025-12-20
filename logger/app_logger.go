package logger

import (
	"flag"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger interface {
	GetLogger(prefix string) Logger
	InitFlags()
	Activate() error
	Stop() error
}

type appLogger struct {
	level  string
	logger *zap.Logger
}

func NewAppLogger() AppLogger {
	logger, _ := zap.NewProduction()
	return &appLogger{
		level:  "info",
		logger: logger,
	}
}

func (a *appLogger) InitFlags() {
	flag.StringVar(
		&a.level,
		"log-level",
		"info",
		"Log level: debug | info | warn | error",
	)
}

func (a *appLogger) Activate() error {
	cfg := zap.NewDevelopmentConfig()
	if a.level == "info" {
		cfg = zap.NewProductionConfig()
	}

	lv, err := zapcore.ParseLevel(a.level)
	if err != nil {
		return err
	}

	cfg.Level = zap.NewAtomicLevelAt(lv)

	logger, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return err
	}

	a.logger = logger
	return nil
}

func (a *appLogger) GetLogger(prefix string) Logger {
	return &zapLogger{
		sugar: a.logger.
			With(zap.String("prefix", prefix)).
			Sugar(),
	}
}

func (a *appLogger) Stop() error {
	if a.logger != nil {
		_ = a.logger.Sync()
	}
	return nil
}
