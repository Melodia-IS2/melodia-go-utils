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
	LayerApp        Layer = "APP"
)

type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     Level     `json:"level"`
	Layer     Layer     `json:"layer"`
	Message   string    `json:"message"`
	Data      any       `json:"data"`
}

type Log struct {
	ID        string `json:"id"`
	AppName   string `json:"app_name"`
	Endpoint  string `json:"endpoint"`
	Method    string `json:"method"`
	RemoteIP  string `json:"remote_ip"`
	UserAgent string `json:"user_agent"`

	Status   int           `json:"status"`
	Bytes    int           `json:"bytes"`
	Duration time.Duration `json:"duration"`

	Timestamp time.Time `json:"timestamp"`
	Entries   []Entry   `json:"entries"`
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
