package middleware

import (
	"base-service/internal/logger"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// LoggingMiddleware логирует каждый запрос
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		requestID := getRequestID(r)

		// загрузка логгера в контекст
		loggerCtx := slog.Default().With(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
		loggerCtx.Info("→ incoming request")

		ctx := logger.LoggerToContext(r.Context(), loggerCtx)
		r = r.WithContext(ctx)

		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lw, r)

		loggerCtx.Debug("← outgoing response",
			"status", lw.statusCode,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// Обёртка над ResponseWriter для получения статуса ответа
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader перехватывает вызов от ResponseWriter
func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// Формирует ID запроса
func getRequestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-Id"); id != "" {
		return id
	}
	return uuid.New().String()
}

func PrintRoutes(r *mux.Router) {
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, errTpl := route.GetPathTemplate()
		if errTpl != nil {
			tpl = "???"
		}

		meths, errMeth := route.GetMethods()
		if errMeth != nil {
			meths = []string{"*"}
		}

		slog.Default().Debug("registered route",
			"method", strings.Join(meths, ", "),
			"path", tpl,
		)
		return nil
	})

	if err != nil {
		slog.Default().Error("failed to walk routes", "error", err)
	}
}
