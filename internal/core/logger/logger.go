package logger

import (
	"go.uber.org/zap"
)

// Log будет доступен всему коду как синглтон.
// var Log *zap.Logger = zap.NewNop()
var log *zap.Logger
var sugar *zap.SugaredLogger

type Logger struct {
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
}

func NewLogger(level string) (*Logger, error) {
	err := initialize(level)
	if err != nil {
		return nil, err
	}

	return &Logger{Log: log, Sugar: sugar}, nil
}

func initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	cfg.Encoding = "console"
	cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	log = zl
	sugar = log.Sugar()
	return nil
}
