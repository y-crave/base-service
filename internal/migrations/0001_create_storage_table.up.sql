-- ~/go/bin/migrate create -ext sql -dir migrations -seq -digits 4 create_storage_table

CREATE TABLE storages
(
    id            UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    bucket_name   VARCHAR(255) NOT NULL,
    service_owner VARCHAR(50)  NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at    TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

COMMENT ON TABLE storage IS 'Хранилище информации о buckets из Minio/S3';

COMMENT ON COLUMN storage.id IS 'Уникальный идентификатор хранилища';
COMMENT ON COLUMN storage.bucket_name IS 'Название бакета в Minio/S3';
COMMENT ON COLUMN storage.service_owner IS 'Владелец/сервис, которому принадлежит бакет';
COMMENT ON COLUMN storage.created_at IS 'Дата и время создания записи';
COMMENT ON COLUMN storage.updated_at IS 'Дата и время последнего обновления';
COMMENT ON COLUMN storage.deleted_at IS 'Мягкое удаление, NULL если запись активна';
