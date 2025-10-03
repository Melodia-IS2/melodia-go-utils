package minio

import (
	"context"
	"io"
)

type MinioBucket interface {
	UploadFile(ctx context.Context, fileName string, file io.Reader) error
	DownloadFile(ctx context.Context, fileName string) (io.Reader, error)
	DeleteFile(ctx context.Context, fileName string) error
}

// Save(context.Context, string, *multipart.FileHeader) error
// 	Get(context.Context, string) error
// 	GetObjectURL(string) string
