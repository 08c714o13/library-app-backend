-- +goose Up
CREATE TABLE "librarian"
(
    email      varchar(255) PRIMARY KEY,
    name       varchar(255) NOT NULL,
    hash_pass  varchar(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOL DEFAULT FALSE NOT NULL
);


-- +goose Down
DROP TABLE IF EXISTS "librarian";
