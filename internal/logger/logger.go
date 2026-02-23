package logger

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init(production bool) {
	if production {
		Log = zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			Output(zerolog.ConsoleWriter{Out: os.Stderr}).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		Log = zerolog.New(os.Stdout).
			Level(zerolog.DebugLevel).
			Output(zerolog.ConsoleWriter{Out: os.Stderr}).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339
}

func With() zerolog.Context {
	return Log.With()
}

func WithTraceID(c *fiber.Ctx) zerolog.Logger {
	traceID := ""
	if c != nil {
		if id, ok := c.Locals("trace_id").(string); ok {
			traceID = id
		}
	}
	return Log.With().Str("trace_id", traceID).Logger()
}

func Info() *zerolog.Event {
	return Log.Info()
}

func Debug() *zerolog.Event {
	return Log.Debug()
}

func Warn() *zerolog.Event {
	return Log.Warn()
}

func Error() *zerolog.Event {
	return Log.Error()
}

func Fatal() *zerolog.Event {
	return Log.Fatal()
}
