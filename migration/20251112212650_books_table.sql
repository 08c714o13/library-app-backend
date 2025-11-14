-- +goose Up
CREATE TABLE "book"
(
    id         uuid PRIMARY KEY,
    title      varchar(255) NOT NULL,
    author     varchar(255) NOT NULL,
    genre      varchar(255) NOT NULL,
    bookcase   varchar(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    user_email varchar(255),
    CONSTRAINT fk_user
        FOREIGN KEY (user_email)
        REFERENCES "user" (email)
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS "book";
