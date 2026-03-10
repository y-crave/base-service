// Package errors предоставляет кастомный тип ошибки TraceError с сохранением стека вызовов
// и метаданных для структурированного логирования (slog).
package errors

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	missingValue = "!MISSING_VALUE"
	badKey       = "!BAD_KEY"
)

// TraceError — ошибка с сообщением, вложенной причиной, опциональным стеком и полями для логов.
type TraceError struct {
	msg      string         // сообщение на этом уровне цепочки
	cause    error          // вложенная ошибка (причина)
	stack    []uintptr      // program counters стека вызовов; заполняется только у первой в цепочке TraceError
	metaData map[string]any // произвольные пары ключ-значение для логирования
}

// NewTraceError создаёт TraceError с сообщением msg, вызываемой ошибкой err и дополнительными полями для логгирования
// в виде произвольных пар ключ–значение keysAndValues, где ключ должен поддерживать string()/Stringer.
// Стек сохраняется только если err == nil или err не является *TraceError.
// При оборачивании уже существующей *TraceError стек не дублируется.
func NewTraceError(msg string, err error, keysAndValues ...any) *TraceError {
	e := TraceError{
		msg:   msg,
		cause: err,
	}

	var te *TraceError
	if err == nil || !errors.As(err, &te) {
		e.stack = captureStack(3)
	}

	e.addMetaData(keysAndValues...)

	return &e
}

// Error возвращает текстовое представление ошибки в виде "msg: cause".
func (e *TraceError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.cause)
	}
	return e.msg
}

// Unwrap возвращает вложенную ошибку; нужен для errors.Is и errors.As.
func (e *TraceError) Unwrap() error {
	return e.cause
}

// LogValue реализует slog.LogValuer: при передаче TraceError в slog выводятся msg, cause_error, stack и metaData.
func (e *TraceError) LogValue() slog.Value {
	attrs := []slog.Attr{slog.String("msg", e.msg)}

	var te *TraceError
	if e.cause != nil {
		if errors.As(e.cause, &te) {
			attrs = append(attrs, slog.Any("cause_error", te)) // рекурсивно вызывается LogValue для вложенной ошибки
		} else {
			attrs = append(attrs, slog.String("cause_error", e.cause.Error()))
		}
	}

	if e.cause == nil || !errors.As(e.cause, &te) {
		stackLines := strings.Split(e.formatStack(), "\n")
		attrs = append(attrs, slog.Any("stack", stackLines))
	}

	for k, v := range e.metaData {
		attrs = append(attrs, slog.Any(k, v))
	}

	return slog.GroupValue(attrs...)
}

// addMetaData добавляет дополнительные поля для логирования, в виде пар {key, value}
func (e *TraceError) addMetaData(keyAndValues ...any) {
	if e.metaData == nil {
		e.metaData = make(map[string]any)
	}

	if len(keyAndValues)%2 != 0 {
		keyAndValues = append(keyAndValues, missingValue)
	}

	for i := 0; i < len(keyAndValues)-1; i += 2 {
		var key string
		switch k := keyAndValues[i].(type) {
		case string:
			key = k
		case fmt.Stringer:
			key = k.String()
		default:
			key = fmt.Sprintf("%s%v", badKey, keyAndValues[i])
		}
		e.metaData[key] = keyAndValues[i+1]
	}
}

// formatStack превращает срез program counters (PC) в читаемые строки
func (e *TraceError) formatStack() string {
	if len(e.stack) == 0 {
		return ""
	}

	frames := runtime.CallersFrames(e.stack) // преобразование массива program counters (PC) в читаемые фреймы стека
	var sb strings.Builder

	for {
		frame, more := frames.Next()

		if !isStdLib(frame.File) {
			sb.WriteString(fmt.Sprintf("%s:%d in %s\n", frame.File, frame.Line, frame.Function))
		}

		if !more {
			break
		}
	}
	return strings.TrimSpace(sb.String())
}

// captureStack сохраняет текущий call stack в виде среза program counters (PC).
// skip — сколько верхних кадров пропустить
func captureStack(skip int) []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[:n]
}

// isStdLib определяет, относится ли путь к файлу к стандартной библиотеке Go
func isStdLib(path string) bool {
	goroot := runtime.GOROOT()

	pathClean := filepath.Clean(path)
	gorootClean := filepath.Clean(goroot)
	rel, err := filepath.Rel(gorootClean, pathClean)

	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, "..")
}
