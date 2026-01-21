CREATE TABLE files
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    path       VARCHAR(255)  NOT NULL,
    size       BIGINT                   DEFAULT 0,
    mime_type  VARCHAR(50)  NOT NULL,
    is_temp    BOOL                     DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

COMMENT
ON TABLE file IS 'Хранилище информации о файлах';

COMMENT ON COLUMN file.id IS 'Уникальный идентификатор файла';
COMMENT ON COLUMN file.name IS 'Оригинальное имя файла';
COMMENT ON COLUMN file.path IS 'Путь к файлу в хранилище Minio/S3';
COMMENT ON COLUMN file.size IS 'Размер файла в байтах';
COMMENT ON COLUMN file.mime_type IS 'MIME-тип файла';
COMMENT ON COLUMN file.is_temp IS 'Является ли файл временным (True - временный, False - постоянный, загружен в Minio/S3)';
COMMENT ON COLUMN file.created_at IS 'Дата и время создания записи в БД';
COMMENT ON COLUMN file.updated_at IS 'Дата и время последнего обновления записи';
COMMENT ON COLUMN file.deleted_at IS 'Мягкое удаление, NULL если запись активна';
