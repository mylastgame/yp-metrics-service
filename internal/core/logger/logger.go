package logger

import (
	"go.uber.org/zap"
)

// Log будет доступен всему коду как синглтон.
// var Log *zap.Logger = zap.NewNop()
var Log *zap.Logger
var Sugar *zap.SugaredLogger

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
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
	Log = zl
	Sugar = Log.Sugar()
	return nil
}
