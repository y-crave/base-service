package contract

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
)

type IFileManager interface {
	UploadFile(context.Context, multipart.File, *multipart.FileHeader) error
	DownloadFile(ctx context.Context, fileName string) (io.Reader, error)
	DeleteFile(ctx context.Context, fileName string) error
}

type IRepo interface {
	AddBucket(BucketName, ServiceOwner string) uuid.UUID
	DeleteBucketByName(BucketName string) error
	DeleteBucketByID(id uuid.UUID) error

	AddFile(BucketID uuid.UUID, Name, Path string, Size int64, MimeType string, IsTemp bool) error
	DeleteFileByID(uuid.UUID) error
}
