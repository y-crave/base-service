package user_repo

import (
	"base-service/internal/domain"
)

func StorageToDomain(storage *Storage) domain.Storage {
	return domain.Storage{
		ID:           storage.ID,
		BucketName:   storage.BucketName,
		ServiceOwner: storage.ServiceOwner,
		CreatedAt:    storage.CreatedAt,
		UpdatedAt:    storage.UpdatedAt,
		DeletedAt:    storage.DeletedAt,
	}
}

func StorageToRepo(storage *domain.Storage) Storage {
	return Storage{
		ID:           storage.ID,
		BucketName:   storage.BucketName,
		ServiceOwner: storage.ServiceOwner,
		CreatedAt:    storage.CreatedAt,
		UpdatedAt:    storage.UpdatedAt,
		DeletedAt:    storage.DeletedAt,
	}
}

func FileToDomain(file *File) domain.File {
	return domain.File{
		ID:        file.ID,
		Name:      file.Name,
		Path:      file.Path,
		Size:      file.Size,
		MimeType:  file.MimeType,
		IsTemp:    file.IsTemp,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
		DeletedAt: file.DeletedAt,
	}
}

func FileToRepo(file *domain.File) File {
	return File{
		ID:        file.ID,
		Name:      file.Name,
		Path:      file.Path,
		Size:      file.Size,
		MimeType:  file.MimeType,
		IsTemp:    file.IsTemp,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
		DeletedAt: file.DeletedAt,
	}
}
