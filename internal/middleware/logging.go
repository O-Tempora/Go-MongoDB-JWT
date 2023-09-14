package middleware

import (
	"context"
	"net/http"

	"golang.org/x/exp/slog"
)

func LogRequest(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.LogAttrs(
				context.Background(),
				slog.LevelInfo,
				"Request: ",
				slog.String("URL", r.URL.Path),
				slog.String("Method", r.Method),
				slog.String("Host", r.Host),
			)
			next.ServeHTTP(w, r)
		})
	}
}
