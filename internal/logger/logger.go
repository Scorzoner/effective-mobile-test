package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapService interface {
	Debug(fields ...interface{})
	Info(fields ...interface{})
	Error(fields ...interface{})
	Fatal(fields ...interface{})
}

type ZapStorage struct {
	*zap.Logger
}

var Zap ZapService = &ZapStorage{zap.NewNop()}

func Init() error {
	config := zap.NewDevelopmentConfig()

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	Zap = &ZapStorage{logger}

	return nil
}

func (z *ZapStorage) Debug(fields ...interface{}) {
	z.Logger.Sugar().Debugln(fields...)
}

func (z *ZapStorage) Info(fields ...interface{}) {
	z.Logger.Sugar().Infoln(fields...)
}

func (z *ZapStorage) Error(fields ...interface{}) {
	z.Logger.Sugar().Errorln(fields...)
}

func (z *ZapStorage) Fatal(fields ...interface{}) {
	z.Logger.Sugar().Fatalln(fields...)
}
