package sql

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/tracelog"
)

type Logger struct {
	logger   *slog.Logger
	omitArgs bool
}

func NewLogger(l *slog.Logger, omitArgs bool) *Logger {
	return &Logger{logger: l, omitArgs: omitArgs}
}

func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		if k == "args" && l.omitArgs {
			continue
		}
		attrs = append(attrs, slog.Any(k, v))
	}

	var lvl slog.Level
	switch level {
	case tracelog.LogLevelTrace:
		lvl = slog.LevelDebug - 1
		attrs = append(attrs, slog.Any("pgx_log_level", level))
	case tracelog.LogLevelDebug:
		lvl = slog.LevelDebug
	case tracelog.LogLevelInfo:
		lvl = slog.LevelInfo
	case tracelog.LogLevelWarn:
		lvl = slog.LevelWarn
	case tracelog.LogLevelError:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelError
		attrs = append(attrs, slog.Any("invalid_pgx_log_level", level))
	}
	l.logger.LogAttrs(ctx, lvl, msg, attrs...)
}
