package minio

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type ObjectData struct {
	Reader       io.Reader
	ContentType  string
	FileName     string
	Size         int64
	UserMetadata map[string]string
}

type MinioBucket interface {
	UploadFile(ctx context.Context, fileName string, file io.Reader, fileSize int64, opts minio.PutObjectOptions) error
	DownloadFile(ctx context.Context, fileName string) (ObjectData, error)
	DeleteFile(ctx context.Context, fileName string) error

	UploadFileHeader(ctx context.Context, fileName string, fileHeader *multipart.FileHeader, opts minio.PutObjectOptions) error
	GetObjectURL(objectName string) string
}
