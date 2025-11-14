-- +goose Up
CREATE TABLE "user"
(
    email      varchar(255) PRIMARY KEY , -- todo number
    is_deleted BOOL DEFAULT FALSE NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS "user";
