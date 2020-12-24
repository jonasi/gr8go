package gr8log

import (
	"context"
)

type ctxKey struct{}

// L returns the logger found on ctx
func L(ctx context.Context) Logger {
	if ctx == nil {
		return NullLogger
	}

	l, ok := ctx.Value(ctxKey{}).(Logger)
	if !ok {
		return NullLogger
	}

	return l
}

// WithLogger stores a new logger on the context
func WithLogger(ctx context.Context, l Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, ctxKey{}, l)
}

// WithArg stores a new logger on the context with k=>v set on it's
// kv
func WithArg(ctx context.Context, k string, v interface{}) context.Context {
	return WithArgs(ctx, Args{k: v})
}

// WithArgs stores a new logger on the context with k=>v set on it's
// kv
func WithArgs(ctx context.Context, args Args) context.Context {
	return WithLogger(ctx, L(ctx).WithArgs(args))
}

// Debug logs with level=debug
func Debug(ctx context.Context, msg string, args ...Args) {
	L(ctx).Log(NewEntry(LevelDebug, msg, args...))
}

// Info logs with level=info
func Info(ctx context.Context, msg string, args ...Args) {
	L(ctx).Log(NewEntry(LevelInfo, msg, args...))
}

// Warn logs with level=warn
func Warn(ctx context.Context, msg string, args ...Args) {
	L(ctx).Log(NewEntry(LevelWarn, msg, args...))
}

// Error logs with level=error
func Error(ctx context.Context, msg string, args ...Args) {
	L(ctx).Log(NewEntry(LevelError, msg, args...))
}

// Fatal logs with level=fatal
func Fatal(ctx context.Context, msg string, args ...Args) {
	L(ctx).Log(NewEntry(LevelFatal, msg, args...))
}
