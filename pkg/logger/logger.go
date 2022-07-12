package logger

import (
	"io"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LOGGER_FILE_STDOUT = "stdout"
	LOGGER_FILE_STDERR = "stderr"

	LOGGER_ENCODER_JSON    = "json"
	LOGGER_ENCODER_CONSOLE = "console"
)

func init() {
	pool = NewPool()
	SetGlobal(Config{})
}

// NewLogger Create a logger, reutrn zap logger.
func NewLogger(cfg Config) *zap.Logger {
	econfig := zapcore.EncoderConfig{
		// FunctionKey:    "function",
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",
		MessageKey:     "message",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// build config
	cfg.Build()

	// log level
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalln(err)
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	// log writer multi with '&'
	// e.g: stdout&/foo/bar.log
	wss := make([]zapcore.WriteSyncer, 0)
	for _, filename := range cfg.Filenames {
		var writer io.Writer
		// if log to console, using color level
		switch filename {
		case LOGGER_FILE_STDOUT:
			writer = os.Stdout
		case LOGGER_FILE_STDERR:
			writer = os.Stderr
		default:
			// log rotate
			writer = &lumberjack.Logger{
				Filename:   filename,
				MaxSize:    cfg.MaxSize,
				MaxAge:     cfg.MaxAge,
				MaxBackups: cfg.MaxBackups,
				LocalTime:  cfg.LocalTime,
				Compress:   cfg.Compress,
			}
			econfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
		wss = append(wss, zapcore.AddSync(writer))
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(wss...)

	// log encoder
	var encoder zapcore.Encoder
	switch cfg.Encoder {
	case LOGGER_ENCODER_CONSOLE:
		encoder = zapcore.NewConsoleEncoder(econfig)
	case LOGGER_ENCODER_JSON:
		encoder = zapcore.NewJSONEncoder(econfig)
	}

	// build logger
	zcore := zapcore.NewCore(encoder, writeSyncer, atomicLevel)
	logger := zap.New(zcore, zap.WithCaller(true))
	return logger
}

// SetGlobal configure a global logger to use,
// after configured, use zap.L() or zap.S(),
// do not need to import package logger again.
func SetGlobal(config Config) {
	zap.ReplaceGlobals(NewLogger(config))
}

// L global logger with field
func L() *zap.Logger {
	return zap.L()
}

// S() global logger without field
func S() *zap.SugaredLogger {
	return zap.S()
}
