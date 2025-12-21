package zaplogger

import (
	"flag"

	"github.com/DatLe328/service-context/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type appLogger struct {
	level       string
	atomicLevel zap.AtomicLevel
	logger      *zap.Logger
}

func NewZapLogger() logger.AppLogger {
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	zl, _ := zap.NewProduction(
		zap.IncreaseLevel(atomicLevel),
	)

	return &appLogger{
		level:       "info",
		atomicLevel: atomicLevel,
		logger:      zl,
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

	a.atomicLevel.SetLevel(lv)
	cfg.Level = a.atomicLevel

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

func (a *appLogger) GetLogger(prefix string) logger.Logger {
	return &zapLogger{
		sugar: a.logger.
			With(zap.String("prefix", prefix)).
			Sugar(),
	}
}

func (a *appLogger) GetLevel() string {
	return a.atomicLevel.Level().String()
}

func (a *appLogger) Stop() error {
	if a.logger != nil {
		_ = a.logger.Sync()
	}
	return nil
}
