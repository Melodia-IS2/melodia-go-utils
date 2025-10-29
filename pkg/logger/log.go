package logger

import (
	"context"
	"time"
)

type Level string

const (
	Debug Level = "DEBUG"
	Audit Level = "AUDIT"
	Error Level = "ERROR"
)

type Layer string

const (
	LayerHandler    Layer = "HANDLER"
	LayerService    Layer = "SERVICE"
	LayerRepository Layer = "REPOSITORY"
)

type Entry struct {
	Timestamp time.Time
	Level     Level
	Layer     Layer
	Message   string
	Data      any
}

type Log struct {
	ID        string
	AppName   string
	Endpoint  string
	Method    string
	RemoteIP  string
	UserAgent string

	Status   int
	Bytes    int
	Duration time.Duration

	Timestamp time.Time
	Entries   []Entry
}

type ctxKey string

const logKey ctxKey = "request_log"

func WithLog(ctx context.Context, l *Log) context.Context {
	return context.WithValue(ctx, logKey, l)
}

func FromContext(ctx context.Context) *Log {
	if l, ok := ctx.Value(logKey).(*Log); ok {
		return l
	}
	return nil
}

func Add(ctx context.Context, level Level, layer Layer, msg string, data any) {
	if l := FromContext(ctx); l != nil {
		l.Entries = append(l.Entries, Entry{
			Timestamp: time.Now(),
			Level:     level,
			Layer:     layer,
			Message:   msg,
			Data:      data,
		})
	}
}
