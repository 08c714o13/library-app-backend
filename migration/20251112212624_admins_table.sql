-- +goose Up
CREATE TABLE "admin"
(
    email     varchar(255) PRIMARY KEY,
    hash_pass varchar(255)
);

-- +goose Down
DROP TABLE IF EXISTS "admin";
