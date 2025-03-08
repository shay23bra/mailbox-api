package logger

import (
	"log/syslog"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger() *Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	useSyslog := os.Getenv("USE_SYSLOG") == "true"
	if useSyslog {
		syslogWriter, err := syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "mailbox-api")
		if err == nil {
			log.Logger = zerolog.New(zerolog.SyslogLevelWriter(syslogWriter))
			return &Logger{logger: log.Logger}
		}
		log.Warn().Err(err).Msg("Failed to connect to syslog, falling back to console logging")
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}

	return &Logger{logger: logger}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	event := l.logger.Debug().Str("level", "debug")
	appendKeyValues(event, args...)
	event.Msg(msg)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	event := l.logger.Info().Str("level", "info")
	appendKeyValues(event, args...)
	event.Msg(msg)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	event := l.logger.Warn().Str("level", "warn")
	appendKeyValues(event, args...)
	event.Msg(msg)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	event := l.logger.Error().Str("level", "error")
	appendKeyValues(event, args...)
	event.Msg(msg)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	event := l.logger.Fatal().Str("level", "fatal")
	appendKeyValues(event, args...)
	event.Msg(msg)
}

func appendKeyValues(event *zerolog.Event, args ...interface{}) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if !ok {
				key = "unknown"
			}
			event.Interface(key, args[i+1])
		}
	}
}
