package log

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
	"github.com/primalskill/golog"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

func MustSetup(env string) *slog.Logger {
	switch env {
	case EnvDev:
		return Dev()
	case EnvProd:
		return Prod()
	default:
		panic("unknown env: " + env)
	}
}

func Dev() *slog.Logger {
	slogOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	opts := &devslog.Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		NewLineAfterLog:   true,
	}

	return slog.New(devslog.NewHandler(os.Stdout, opts))
}

func Prod() *slog.Logger {
	return golog.NewProduction()
}
