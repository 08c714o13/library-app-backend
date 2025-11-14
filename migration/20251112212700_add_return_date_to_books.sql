-- +goose Up
ALTER TABLE "book" ADD COLUMN return_date timestamp;

-- +goose Down
ALTER TABLE "book" DROP COLUMN return_date;
