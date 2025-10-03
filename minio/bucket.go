package minio

import (
	"context"
	"errors"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioBucketImpl struct {
	client     *minio.Client
	bucketName string
}

func NewMinioBucket(client *minio.Client, bucketName string) MinioBucket {
	return &MinioBucketImpl{
		client:     client,
		bucketName: bucketName,
	}
}

func (b *MinioBucketImpl) UploadFile(ctx context.Context, fileName string, file io.Reader) error {
	return errors.New("not implemented")
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
