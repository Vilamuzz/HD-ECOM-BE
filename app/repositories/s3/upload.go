package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func EnsureBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}

// UploadFile uploads a file to MinIO and returns the object name
func (r *s3Repository) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("attachments/%s_%s%s", time.Now().Format("20060102"), uuid.New().String(), ext)

	// Upload to MinIO
	_, err = r.client.PutObject(ctx, r.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object %s to bucket %s: %w", objectName, r.bucketName, err)
	}

	return objectName, nil
}

// DeleteFile deletes a file from MinIO
func (r *s3Repository) DeleteFile(ctx context.Context, objectName string) error {
	return r.client.RemoveObject(ctx, r.bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetFileURL generates a presigned URL for downloading a file
func (r *s3Repository) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	u, err := r.client.PresignedGetObject(ctx, r.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// DeleteFile deletes a file from MinIO (standalone function for backward compatibility)
func DeleteFile(ctx context.Context, client *minio.Client, bucketName, objectName string) error {
	return client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetFileURL generates a presigned URL for downloading a file (standalone function for backward compatibility)
func GetFileURL(ctx context.Context, client *minio.Client, bucketName, objectName string, expiry time.Duration) (string, error) {
	u, err := client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
