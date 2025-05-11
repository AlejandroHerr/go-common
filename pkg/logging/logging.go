package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/golang-cz/devslog"
)

type (
	config struct {
		Environment string     // "development", "production", "staging"
		Level       slog.Level // "debug", "info", "warn", "error"
		App         string
		Version     string
		Commit      string
		BuildTime   string
		GoVersion   string
		CtxKeys     map[any]string // Maps context keys to attribute names
	}
	ConfigFunc func(*config)
)

func NewLogger(cfgs ...ConfigFunc) *slog.Logger {
	cfg := &config{
		Environment: "development",
		Level:       slog.LevelDebug,
		App:         "n/a",
		Version:     "n/a",
		Commit:      "n/a",
		BuildTime:   "n/a",
		GoVersion:   "n/a",
		CtxKeys:     make(map[any]string),
	}

	for _, c := range cfgs {
		c(cfg)
	}

	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.Environment == "development",
	}

	var handler slog.Handler

	switch strings.ToLower(cfg.Environment) {
	case "development":
		// Use pretty text handler for development
		handler = devslog.NewHandler(os.Stdout, &devslog.Options{ //nolint:exhaustruct //no need for all
			HandlerOptions: opts,
		})
	default:
		// Use JSON for production/staging
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	if len(cfg.CtxKeys) > 0 {
		handler = NewContextHandler(handler, cfg.CtxKeys)
	}

	// Create base logger with common attributes
	logger := slog.New(handler).With(
		"app", cfg.App,
		"environment", cfg.Environment,
		"version", cfg.Version,
		"commit", cfg.Commit,
		"build_time", cfg.BuildTime,
		"go_version", cfg.GoVersion,
	)

	return logger
}

func WithLevel(lvl string) ConfigFunc {
	return func(cfg *config) {
		slevel := levelFromString(lvl, slog.LevelDebug)
		cfg.Level = slevel
	}
}

func WithEnvironment(env string) ConfigFunc {
	return func(cfg *config) {
		cfg.Environment = env
	}
}

func WithApp(app string) ConfigFunc {
	return func(cfg *config) {
		cfg.App = strings.ToLower(app)
	}
}

func WithVersion(version string) ConfigFunc {
	return func(cfg *config) {
		cfg.Version = version
	}
}

func WithCommit(commit string) ConfigFunc {
	return func(cfg *config) {
		cfg.Commit = commit
	}
}

func WithBuildTime(buildTime string) ConfigFunc {
	return func(cfg *config) {
		cfg.BuildTime = buildTime
	}
}

func WithGoVersion(goVersion string) ConfigFunc {
	return func(cfg *config) {
		cfg.GoVersion = goVersion
	}
}

func WithCtxKeys(keys map[any]string) ConfigFunc {
	return func(cfg *config) {
		cfg.CtxKeys = keys
	}
}

func levelFromString(level string, defaultLevel slog.Level) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return defaultLevel
	}
}
