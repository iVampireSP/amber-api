package logger

import "go.uber.org/zap"

type Logger struct {
	Logger *zap.SugaredLogger
}

func NewZapLogger() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
		return nil
	}
	return &Logger{Logger: logger.Sugar()}
}
