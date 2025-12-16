package app

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

// PDOKMetadataToolCLI contains the logic of the PDOK Metadata Tool CLI.
var PDOKMetadataToolCLI = &cli.Command{
	Name:                  "pmt",
	Usage:                 "PDOK Metadata Tool - This tool is set up to handle various metadata related tasks.",
	EnableShellCompletion: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "log-level",
			Usage: "Set log level: debug, info, warn, error (env: PMT_LOG_LEVEL)",
			Value: "info",
		},
	},
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		// Allow env var override
		val := cmd.String("log-level")
		if env := os.Getenv("PMT_LOG_LEVEL"); strings.TrimSpace(env) != "" {
			val = env
		}
		lvl := parseLogLevel(val)
		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl})
		slog.SetDefault(slog.New(handler))

		return ctx, nil
	},
}

func parseLogLevel(s string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
