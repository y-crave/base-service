package storage_repo

import "time"

// Storage представляет модель таблицы storage в базе данных
type Storage struct {
	ID           string     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`      // Уникальный идентификатор хранилища
	BucketName   string     `gorm:"column:bucket_name;type:varchar(255);not null"`                 // Название бакета в Minio/S3
	ServiceOwner string     `gorm:"column:service_owner;type:varchar(50);not null"`                // Владелец/сервис, которому принадлежит бакет
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp with time zone;default:now()"` // Дата и время создания записи
	UpdatedAt    time.Time  `gorm:"column:created_at;type:timestamp with time zone;default:now()"` // Дата и время последнего обновления
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:timestamp with time zone;default:null"`  // Мягкое удаление, NULL если запись активна
}

func (Storage) TableName() string {
	return "storages"
}

// File представляет модель таблицы file в базе данных
type File struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`      // Уникальный идентификатор файла
	Name      string     `gorm:"column:name;type:varchar(255);not null"`                        // Оригинальное имя файла
	Path      string     `gorm:"column:path;type:varchar(255);not null"`                        // Путь к файлу в хранилище Minio/S3
	Size      int64      `gorm:"column:size;type:bigint;default:0"`                             // Размер файла в байтах
	MimeType  string     `gorm:"column:mime_type;type:varchar(50);not null"`                    // MIME-тип файла
	IsTemp    bool       `gorm:"column:is_temp;type:bool;default:true"`                         // Является ли файл временным (True - временный, False - постоянный, загружен в Minio/S3)
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp with time zone;default:now()"` // Дата и время создания записи в БД
	UpdatedAt time.Time  `gorm:"column:created_at;type:timestamp with time zone;default:now()"` // Дата и время последнего обновления записи
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp with time zone;default:null"`  // Мягкое удаление, NULL если запись активна
}

func (File) TableName() string {
	return "files"
}
