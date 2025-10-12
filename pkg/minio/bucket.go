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

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		policy := fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": ["*"]},
                    "Action": ["s3:GetObject"],
                    "Resource": ["arn:aws:s3:::%s/*"]
                }
            ]
        }`, bucketName)

		err = client.SetBucketPolicy(context.Background(), bucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	return &MinioBucketImpl{
		client:         client,
		bucketName:     bucketName,
		publicEndpoint: publicEndpoint,
		useSSL:         useSSL,
	}, nil
}

func (b *MinioBucketImpl) UploadFile(ctx context.Context, fileName string, file io.Reader, fileSize int64, opts minio.PutObjectOptions) error {
	_, err := b.client.PutObject(ctx, b.bucketName, fileName, file, fileSize, opts)
	if err != nil {
		return err
	}
	return nil
}

func (b *MinioBucketImpl) DownloadFile(ctx context.Context, fileName string) (ObjectData, error) {
	obj, err := b.client.GetObject(ctx, b.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return ObjectData{}, err
	}

	info, err := obj.Stat()
	if err != nil {
		return ObjectData{}, err
	}

	return ObjectData{
		Reader:       obj,
		ContentType:  info.ContentType,
		FileName:     info.UserMetadata["original_filename"],
		Size:         info.Size,
		UserMetadata: info.UserMetadata,
	}, nil
}

func (b *MinioBucketImpl) DeleteFile(ctx context.Context, fileName string) error {
	return b.client.RemoveObject(ctx, b.bucketName, fileName, minio.RemoveObjectOptions{})
}

func (b *MinioBucketImpl) UploadFileHeader(ctx context.Context, fileName string, fileHeader *multipart.FileHeader, opts minio.PutObjectOptions) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if opts.ContentType == "" {
		opts.ContentType = fileHeader.Header.Get("Content-Type")
	}

	if opts.UserMetadata == nil {
		opts.UserMetadata = map[string]string{}
	}
	opts.UserMetadata["original_filename"] = fileHeader.Filename

	return b.UploadFile(ctx, fileName, file, fileHeader.Size, opts)
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
