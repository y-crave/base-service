package config

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware логирует каждый запрос
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логируем начало запроса
		slog.Info("→ incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		// Оборачиваем ResponseWriter, чтобы получить статус ответа
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}

		// Выполняем следующий обработчик
		next.ServeHTTP(lw, r)

		// Логируем завершение запроса
		duration := time.Since(start)
		slog.Info("← outgoing response",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lw.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
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

		slog.Info("registered route",
			"method", strings.Join(meths, ", "),
			"path", tpl,
		)
		return nil
	})

	if err != nil {
		slog.Error("failed to walk routes", "error", err)
	}
}
