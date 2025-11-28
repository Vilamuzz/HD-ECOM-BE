package helpers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient() *minio.Client {
	// Get from environment variables with defaults
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000" // ✅ Use API port 9000
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:       useSSL,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		log.Fatalf("failed to init minio: %v", err)
	}
	return client
}

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
func UploadFile(ctx context.Context, client *minio.Client, bucketName string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("attachments/%s_%s%s", time.Now().Format("20060102"), uuid.New().String(), ext)

	// Ensure the target bucket exists (create it if not)
	if err := EnsureBucket(ctx, client, bucketName); err != nil {
		log.Printf("[minio] EnsureBucket failed for bucket=%s: %v", bucketName, err)
		return "", fmt.Errorf("bucket %s does not exist and could not be created: %w", bucketName, err)
	}

	// Upload to MinIO
	_, err = client.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		log.Printf("[minio] PutObject failed: bucket=%s object=%s err=%v", bucketName, objectName, err)
		return "", fmt.Errorf("failed to upload object %s to bucket %s: %w", objectName, bucketName, err)
	}

	return objectName, nil
}

// DeleteFile deletes a file from MinIO
func DeleteFile(ctx context.Context, client *minio.Client, bucketName, objectName string) error {
	return client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetFileURL generates a presigned URL for downloading a file
func GetFileURL(ctx context.Context, client *minio.Client, bucketName, objectName string, expiry time.Duration) (string, error) {
	u, err := client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
