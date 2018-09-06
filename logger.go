package tinnitus

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	logger, err := Config.Logger.Build()

	if err != nil {
		panic(err)
	}

	Logger = logger.Sugar()
}
