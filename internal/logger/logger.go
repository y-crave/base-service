// Package logger предоставляет создание логгера на основе slog,
// функции для сохранения/извлечения логгера из context.Context.
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// contextKey — тип ключа для хранения логгера в контексте
type contextKey string

const (
	// LoggerKey — ключ, по которому в context хранится *slog.Logger.
	LoggerKey contextKey = "logger"
)

// NewLogger создаёт логгер с уровнем strLevel ("DEBUG", "INFO", "WARN", "ERROR"),
// выводом в stdout в JSON и кастомным форматом времени и уровня.
func NewLogger(strLevel string) *slog.Logger {

	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(strLevel))); err != nil {
		level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			if a.Key == slog.TimeKey {
				return slog.String("time", a.Value.Time().Format("02.01.2006 15:04:05"))
			}

			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				return slog.String("level", level.String())
			}
			return a
		},
	}))
}

// LoggerToContext сохраняет логгер l в контексте и возвращает новый контекст.
func LoggerToContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, l)
}

// LoggerFromContext возвращает логгер из контекста.
// Если логгер в контексте не задан, возвращается slog.Default().
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
