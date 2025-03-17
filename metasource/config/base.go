package config

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

func SetLogger() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
}
