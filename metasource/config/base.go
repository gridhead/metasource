package config

import (
	"fmt"
	"github.com/lmittmann/tint"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		date := time.Now()
		slog.Log(nil, slog.LevelInfo, fmt.Sprintf("[%s] <%s> %s", date.Format(time.RFC3339), r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func SetLogger() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
}
