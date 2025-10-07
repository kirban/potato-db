package logger

import (
	"fmt"
	"os"

	"github.com/kirban/potato-db/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	lvl, err := parseZapLevel(cfg.App.LogLevel)
	if err != nil {
		fmt.Printf("invalid log level '%s': %v\n", cfg.App.LogLevel, err)
		return nil, err
	}

	var encoding string
	if lvl == zap.DebugLevel {
		encoding = "console"
	} else {
		encoding = "json"
	}

	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		Development:      lvl == zap.DebugLevel,
		Encoding:         encoding,
		OutputPaths:      []string{cfg.App.LogOutput},
		ErrorOutputPaths: []string{cfg.App.LogOutput},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	hook := func(e zapcore.Entry) error {
		switch e.Level {
		case zapcore.PanicLevel:
			formatPanic(e)
		case zapcore.FatalLevel:
			formatFatal(e)
		}
		return nil
	}

	logger, err := zapCfg.Build(
		zap.AddStacktrace(zapcore.PanicLevel),
		zap.AddCaller(),
		zap.Hooks(hook),
	)

	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return nil, err
	}

	return logger, nil
}

func parseZapLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zap.DebugLevel, nil
	case "info":
		return zap.InfoLevel, nil
	case "warn":
		return zap.WarnLevel, nil
	case "error":
		return zap.ErrorLevel, nil
	case "panic":
		return zap.PanicLevel, nil
	case "fatal":
		return zap.FatalLevel, nil
	default:
		return zap.InfoLevel, fmt.Errorf("unknown level: %s", level)
	}
}

func formatPanic(e zapcore.Entry) {
	fmt.Fprintf(os.Stderr, "PANIC: %s | time=%s | caller=%s\n", e.Message, e.Time.Format("2006-01-02T15:04:05Z07:00"), e.Caller.TrimmedPath())
	if e.Stack != "" {
		fmt.Fprintln(os.Stderr, e.Stack)
	}
}

func formatFatal(e zapcore.Entry) {
	fmt.Fprintf(os.Stderr, "FATAL: %s | time=%s | caller=%s\n", e.Message, e.Time.Format("2006-01-02T15:04:05Z07:00"), e.Caller.TrimmedPath())
	if e.Stack != "" {
		fmt.Fprintln(os.Stderr, e.Stack)
	}
}
