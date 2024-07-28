package logger

import "go.uber.org/zap"

type Logger struct {
	Sugar  *zap.SugaredLogger
	Logger *zap.Logger
}

func NewZapLogger() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
		return nil
	}
	return &Logger{Sugar: logger.Sugar(), Logger: logger}
}
