package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init initializes the global logger with optimized production settings.
// Call this once at application startup.
func Init() {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",             // Structured, machine-readable logs
		OutputPaths:      []string{"stdout"}, // Logs go to stdout
		ErrorOutputPaths: []string{"stderr"}, // Errors to stderr
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			TimeKey:     "ts",
			EncodeTime:  zapcore.EpochTimeEncoder,      // Fastest timestamp
			EncodeLevel: zapcore.LowercaseLevelEncoder, // "info", "error", etc.
		},
	}

	var err error
	Log, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}
