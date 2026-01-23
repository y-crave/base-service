package user_repo

import "time"

// Storage представляет модель таблицы storage в базе данных
type Storage struct {
	ID           string     `db:"id"`            // Уникальный идентификатор хранилища
	BucketName   string     `db:"bucket_name"`   // Название бакета в Minio/S3
	ServiceOwner string     `db:"service_owner"` // Владелец/сервис, которому принадлежит бакет
	CreatedAt    time.Time  `db:"created_at"`    // Дата и время создания записи
	UpdatedAt    time.Time  `db:"updated_at"`    // Дата и время последнего обновления
	DeletedAt    *time.Time `db:"deleted_at"`    // Мягкое удаление, NULL если запись активна
}

// File представляет модель таблицы file в базе данных
type File struct {
	ID        string     `db:"id"`         // Уникальный идентификатор файла
	Name      string     `db:"name"`       // Оригинальное имя файла
	Path      string     `db:"path"`       // Путь к файлу в хранилище Minio/S3
	Size      int64      `db:"size"`       // Размер файла в байтах
	MimeType  string     `db:"mime_type"`  // MIME-тип файла
	IsTemp    bool       `db:"is_temp"`    // Является ли файл временным (True - временный, False - постоянный, загружен в Minio/S3)
	CreatedAt time.Time  `db:"created_at"` // Дата и время создания записи в БД
	UpdatedAt time.Time  `db:"updated_at"` // Дата и время последнего обновления записи
	DeletedAt *time.Time `db:"deleted_at"` // Мягкое удаление, NULL если запись активна
}
