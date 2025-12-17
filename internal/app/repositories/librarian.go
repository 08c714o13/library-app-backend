package repositories

import (
	"database/sql"
	"errors"

	"library/internal/app/domain"
)

type LibrarianRepo struct {
	db *sql.DB
}

func NewLibrarianRepo(db *sql.DB) LibrarianRepo {
	return LibrarianRepo{db: db}
}

func (r LibrarianRepo) CreateLibrarian(email, name, hashPassword string) error {
	_, err := r.db.Exec(`
        INSERT INTO "librarian" (email, name, hash_pass, is_deleted, created_at, updated_at) 
        VALUES ($1, $2, $3, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT (email) 
        DO UPDATE SET 
            name = EXCLUDED.name,
            hash_pass = EXCLUDED.hash_pass,
            is_deleted = false,
            updated_at = CURRENT_TIMESTAMP
    `, email, name, hashPassword)
	return err
}

func (r LibrarianRepo) GetLibrarianByEmail(email string) (string, error) {
	var hashPassword string
	err := r.db.QueryRow(`SELECT hash_pass FROM "librarian" WHERE email = $1 AND is_deleted = false`, email).Scan(&hashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("librarian not found")
		}
		return "", err
	}
	return hashPassword, nil
}

func (r LibrarianRepo) GetAllLibrarians() ([]domain.Librarian, error) {
	rows, err := r.db.Query(`SELECT email, name, created_at, updated_at FROM "librarian" WHERE is_deleted = false ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var librarians []domain.Librarian
	for rows.Next() {
		var lib domain.Librarian
		if err := rows.Scan(&lib.Email, &lib.Name, &lib.CreatedAt, &lib.UpdatedAt); err != nil {
			return nil, err
		}
		librarians = append(librarians, lib)
	}
	return librarians, nil
}

func (r LibrarianRepo) DeleteLibrarian(email string) error {
	_, err := r.db.Exec(`UPDATE "librarian" SET is_deleted = true, updated_at = CURRENT_TIMESTAMP WHERE email = $1`, email)
	return err
}

func (r LibrarianRepo) LibrarianExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM "librarian" WHERE email = $1 AND is_deleted = false)`, email).Scan(&exists)
	return exists, err
}
