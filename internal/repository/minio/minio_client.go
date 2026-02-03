package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"storage-service/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient(minio *minio.Client) *MinioClient {
	return &MinioClient{client: minio}
}

func ConnectToMinio(cfg *config.Config) (*minio.Client, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	return client, err
}

func (m *MinioClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return m.client.BucketExists(ctx, bucketName)
}

func (m *MinioClient) CreateBucket(ctx context.Context, bucketName string) error {
	exists, err := m.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if exists {
		return nil
	}

	err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	// доступ для всех и любых действий
	policy := `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": ["*"]},
                    "Action": "s3:*",
                    "Resource": [
                        "arn:aws:s3:::` + bucketName + `",
                        "arn:aws:s3:::` + bucketName + `/*"
                    ]
                }
            ]
        }`

	err = m.client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		log.Printf("Warning: failed to set bucket policy: %s, %v", bucketName, err)
	}

	return nil
}

// удаление только если бакет пустой
func (m *MinioClient) DeleteBucket(ctx context.Context, bucketName string) error {
	return m.client.RemoveBucket(ctx, bucketName)
}

func (m *MinioClient) GetBucketSize(ctx context.Context, bucketName string) (int64, error) {
	var totalSize int64

	objectList := m.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	for object := range objectList {
		if object.Err != nil {
			return 0, fmt.Errorf("failed to list object in bucket %s, %w", bucketName, object.Err)
		}
		totalSize += object.Size
	}
	return totalSize, nil
}

func (m *MinioClient) UploadFileFromReader(ctx context.Context, bucketName, objectName string, reader io.Reader,
	objectSize int64, contentType string) (minio.UploadInfo, error) {

	if ok, err := m.BucketExists(ctx, bucketName); !ok {
		return minio.UploadInfo{}, err
	}
	return m.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})
}

func (m *MinioClient) DownloadFile(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (m *MinioClient) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return m.client.ListBuckets(ctx)
}

func (m *MinioClient) ListNameFiles(ctx context.Context, bucketName string) ([]string, error) {
	var files []string

	objectList := m.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	for object := range objectList {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to get list of files in bucket %s, %w", bucketName, object.Err)
		}
		files = append(files, object.Key)
	}
	return files, nil
}

func (m *MinioClient) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}
