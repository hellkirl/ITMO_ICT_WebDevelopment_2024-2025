package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type ZapLogger struct {
	sugar *zap.Logger
}

var (
	instance *ZapLogger
	once     sync.Once
	err      error
)

func NewZapLogger() (*ZapLogger, error) {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, e := config.Build()
		if e != nil {
			err = e
			return
		}
		instance = &ZapLogger{
			sugar: logger,
		}
	})
	return instance, err
}

func (z *ZapLogger) Debug(msg string, fields ...zap.Field) {
	z.sugar.Debug(msg, fields...)
}

func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.sugar.Info(msg, fields...)
}

func (z *ZapLogger) Warn(msg string, fields ...zap.Field) {
	z.sugar.Warn(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.sugar.Error(msg, fields...)
}

func (z *ZapLogger) Sync() error {
	return z.sugar.Sync()
}
