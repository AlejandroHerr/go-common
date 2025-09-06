package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func RequestLoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status and size
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			requestID, ok := r.Context().Value(chimiddleware.RequestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}

			reqLogger := logger.With(
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("referer", r.Referer()),
				slog.String("host", r.Host),
			)

			if r.URL.RawQuery != "" {
				reqLogger = logger.With(slog.Any("query", r.URL.Query()))
			}

			// Log the route pattern if available from Chi
			if routePattern := chi.RouteContext(r.Context()).RoutePattern(); routePattern != "" {
				reqLogger = logger.With(slog.String("route_pattern", routePattern))
			}

			// Debug level for request start
			reqLogger.Debug("request started")

			next.ServeHTTP(ww, r)

			// Log request completion with outcome information
			status := ww.Status()
			duration := float64(time.Since(start).Microseconds()) / float64(time.Microsecond)

			// Choose log level based on status code
			attrs := []any{
				"status", status,
				"bytes", ww.BytesWritten(),
				"duration_ms", duration,
			}

			// Add content type if present
			if ct := ww.Header().Get("Content-Type"); ct != "" {
				attrs = append(attrs, "content_type", ct)
			}

			switch {
			case status >= http.StatusInternalServerError:
				reqLogger.Error("server error", attrs...)
			case status >= http.StatusBadRequest:
				reqLogger.Warn("client error", attrs...)
			case status >= http.StatusMultipleChoices:
				reqLogger.Info("redirect", attrs...)
			default:
				reqLogger.Info("request completed", attrs...)
			}
		})
	}
}
