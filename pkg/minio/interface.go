package minio

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type MinioBucket interface {
	UploadFile(ctx context.Context, fileName string, file io.Reader, fileSize int64, opts minio.PutObjectOptions) error
	DownloadFile(ctx context.Context, fileName string) (io.Reader, error)
	DeleteFile(ctx context.Context, fileName string) error

	UploadFileHeader(ctx context.Context, fileName string, fileHeader *multipart.FileHeader, opts minio.PutObjectOptions) error
	GetObjectURL(objectName string) string
}
