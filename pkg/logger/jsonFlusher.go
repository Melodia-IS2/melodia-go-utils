package logger

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type JSONFlusher struct {
	writer io.Writer
}

func NewJSONFlusher(path string) *JSONFlusher {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil
		}
	}
	writer, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	return &JSONFlusher{writer: writer}
}

func (f *JSONFlusher) Flush(ctx context.Context, log *Log) error {
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	_, err = f.writer.Write(append(data, '\n'))
	return err
}
