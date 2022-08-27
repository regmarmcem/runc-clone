package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger            *zap.SugaredLogger
	debugOption       bool
	debugOptionStatus *bool
)

func InitLogger(debug bool) (err error) {
	debugOption = debug
	debugOptionStatus = &debugOption
	Logger, err = newLogger(debug)
	if err != nil {
		return fmt.Errorf("error in new logger: %w", err)
	}
	fmt.Println(*debugOptionStatus)
	return nil
}

func newLogger(debug bool) (*zap.SugaredLogger, error) {
	level := zap.NewAtomicLevel()
	if debug {
		level.SetLevel(zapcore.DebugLevel)
	} else {
		level.SetLevel(zapcore.InfoLevel)
	}

	zapConfig := zap.Config{
		Level:             level,
		Development:       debug,
		Encoding:          "console",
		DisableStacktrace: true,
		DisableCaller:     true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap config: %w", err)
	}

	return logger.Sugar(), nil
}

func Fatal(err error) {
	if debugOption {
		Logger.Fatalf("%+v", err)
	}
	Logger.Fatal(err)
}
