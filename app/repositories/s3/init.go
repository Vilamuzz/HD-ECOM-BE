package s3

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type s3Repository struct {
	bucketName string
	client     *minio.Client
}

func NewS3Repository(contextTimeout time.Duration) *s3Repository {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	endpoint := os.Getenv("MINIO_ENDPOINT")
	access := os.Getenv("MINIO_ACCESS_KEY")
	secret := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	// Return nil if MinIO is not configured
	if endpoint == "" {
		log.Println("MinIO endpoint not configured, S3 repository will be disabled")
		return nil
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(access, secret, ""),
		Secure:       useSSL,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		log.Fatalf("failed to init minio: %v", err)
	}

	// determine bucket name (fall back to default used elsewhere)
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "my-bucket"
	}

	// Ensure the bucket exists (uses ctx created above)
	if err := EnsureBucket(ctx, client, bucketName); err != nil {
		log.Printf("[s3] EnsureBucket failed for bucket=%s: %v", bucketName, err)
	}

	return &s3Repository{
		bucketName: bucketName,
		client:     client,
	}
}
