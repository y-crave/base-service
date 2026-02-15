package domain

import "time"

// Storage представляет модель таблицы storage в базе данных
type Storage struct {
	ID           string     // Уникальный идентификатор хранилища
	BucketName   string     // Название бакета в Minio/S3
	ServiceOwner string     // Владелец/сервис, которому принадлежит бакет
	CreatedAt    time.Time  // Дата и время создания записи
	UpdatedAt    time.Time  // Дата и время последнего обновления
	DeletedAt    *time.Time // Мягкое удаление, NULL если запись активна
}

// File представляет модель таблицы file в базе данных
type File struct {
	ID        string     // Уникальный идентификатор файла
	Name      string     // Оригинальное имя файла
	Path      string     // Путь к файлу в хранилище Minio/S3
	Size      int64      // Размер файла в байтах
	MimeType  string     // MIME-тип файла
	IsTemp    bool       // Является ли файл временным (True - временный, False - постоянный, загружен в Minio/S3)
	CreatedAt time.Time  // Дата и время создания записи в БД
	UpdatedAt time.Time  // Дата и время последнего обновления записи
	DeletedAt *time.Time // Мягкое удаление, NULL если запись активна
}
