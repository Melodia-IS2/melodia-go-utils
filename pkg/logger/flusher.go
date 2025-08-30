package logger

import "context"

type Flusher interface {
	Flush(ctx context.Context, log *Log) error
}
