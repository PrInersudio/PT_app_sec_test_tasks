package nulllogger

import (
	"context"
	"log/slog"
)

type NullLogger struct{}

func (h *NullLogger) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (h *NullLogger) Handle(ctx context.Context, r slog.Record) error {
	return nil
}

func (h *NullLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *NullLogger) WithGroup(name string) slog.Handler {
	return h
}
