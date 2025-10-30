package logger

import (
	"context"
	"fmt"
)

type Flusher interface {
	Flush(ctx context.Context, log *Log) error
}

type multiFlusher struct {
	flushers []Flusher
}

func (m *multiFlusher) Flush(ctx context.Context, log *Log) error {
	var errs []error

	for _, f := range m.flushers {
		if err := f.Flush(ctx, log); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiFlusher errors: %v", errs)
	}

	return nil
}

type FlusherBuilder struct {
	flushers []Flusher
}

func NewFlusherBuilder() *FlusherBuilder {
	return &FlusherBuilder{flushers: []Flusher{}}
}

func (b *FlusherBuilder) Add(f Flusher) *FlusherBuilder {
	if f != nil {
		b.flushers = append(b.flushers, f)
	}
	return b
}

func (b *FlusherBuilder) Build() Flusher {
	if len(b.flushers) == 0 {
		return nil
	}
	if len(b.flushers) == 1 {
		return b.flushers[0]
	}
	return &multiFlusher{flushers: b.flushers}
}
