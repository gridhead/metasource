package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		date := time.Now()
		slog.Log(nil, slog.LevelInfo, fmt.Sprintf("[%s] <%s> %s", date.Format(time.RFC3339), r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}
