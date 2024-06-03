package internal

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitializeLogger(level string) (*zap.Logger, error) {
	levelMap := map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
	}

	logLevel, exists := levelMap[level]
	if !exists {
		return nil, fmt.Errorf("invalid log level: %s", level)
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Encoding:         "console",
		EncoderConfig:    getEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return config.Build()
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%s:", level.CapitalString()))
}
