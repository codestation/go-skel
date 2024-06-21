package logger

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

// InitLogger initializes the logger. If the log-format is not specified, it will default to JSON if the output is not a terminal.
// Required viper variables: log-format as string and debug as bool
func InitLogger() {
	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))

	switch viper.GetString("log-format") {
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	case "logfmt":
	case "":
		if !isTerminal {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
		}
	default:
		slog.Error("Invalid log format specified")
		os.Exit(1)
	}

	if viper.GetString("log-format") == "logfmt" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}
