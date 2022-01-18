package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	config zap.Config
	// Sugar performance is nice, but not critical
	Sugar *zap.SugaredLogger
	// Desugar even faster than the SugaredLogger and allocates far less, but it only supports strongly-typed, structured logging.
	Desugar *zap.Logger
)

func init() {
	if os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "test" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// colorful log output
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	defer logger.Sync() // flushes buffer, if any

	Sugar = logger.Sugar()
	Desugar = logger
}
