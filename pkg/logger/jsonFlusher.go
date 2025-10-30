package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type JSONFlusher struct {
	writer io.Writer
}

func NewJSONFlusher(path string) (*JSONFlusher, error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}
	writer, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return &JSONFlusher{writer: writer}, nil
}

func (f *JSONFlusher) Flush(ctx context.Context, log *Log) error {
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	_, err = f.writer.Write(append(data, '\n'))
	return err
}
