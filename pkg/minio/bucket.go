package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type MinioBucketImpl struct {
	client         *minio.Client
	bucketName     string
	publicEndpoint string
	useSSL         bool
}

func NewMinioBucket(client *minio.Client, bucketName string, publicEndpoint string, useSSL bool) (MinioBucket, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if bucketName == "" {
		return nil, errors.New("bucket name is empty")
	}

	return &MinioBucketImpl{
		client:         client,
		bucketName:     bucketName,
		publicEndpoint: publicEndpoint,
		useSSL:         useSSL,
	}, nil
}

func (b *MinioBucketImpl) UploadFile(ctx context.Context, fileName string, file io.Reader, fileSize int64) error {
	_, err := b.client.PutObject(ctx, b.bucketName, fileName, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (b *MinioBucketImpl) DownloadFile(ctx context.Context, fileName string) (io.Reader, error) {
	obj, err := b.client.GetObject(ctx, b.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (b *MinioBucketImpl) DeleteFile(ctx context.Context, fileName string) error {
	return b.client.RemoveObject(ctx, b.bucketName, fileName, minio.RemoveObjectOptions{})
}

func (b *MinioBucketImpl) UploadFileHeader(ctx context.Context, fileName string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	return b.UploadFile(ctx, fileName, file, fileHeader.Size)
}

func (b *MinioBucketImpl) GetObjectURL(objectName string) string {
	if objectName == "" {
		return ""
	}

	protocol := "http"
	if b.useSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, b.publicEndpoint, b.bucketName, objectName)
}
