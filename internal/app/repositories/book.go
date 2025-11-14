package repositories

import (
	"database/sql"
	"time"

	"library/internal/app/domain"

	"github.com/google/uuid"
)

type BookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) BookRepo {
	return BookRepo{db: db}
}

func (r BookRepo) CreateBook(book domain.CreateBookPayload) (string, error) {
	id := uuid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO "book" (id, title, author, genre, bookcase) VALUES ($1, $2, $3, $4, $5)`,
		id, book.Title, book.Author, book.Genre, book.Bookcase,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r BookRepo) GetAllBooks() ([]domain.Book, error) {
	rows, err := r.db.Query(`
		SELECT id, title, author, genre, bookcase, user_email, return_date, created_at, updated_at
		FROM "book"
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Bookcase,
			&book.UserEmail, &book.ReturnDate, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r BookRepo) GetAllBooksWithUsers() ([]domain.BookWithUser, error) {
	rows, err := r.db.Query(`
		SELECT
			b.id, b.title, b.author, b.genre, b.bookcase,
			b.return_date, b.created_at, b.updated_at,
			u.email, u.created_at
		FROM "book" b
		LEFT JOIN "user" u ON b.user_email = u.email AND u.is_deleted = false
		ORDER BY b.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var books []domain.BookWithUser
	for rows.Next() {
		var book domain.BookWithUser
		var userEmail *string
		var userCreatedAt *time.Time

		if err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.Genre, &book.Bookcase,
			&book.ReturnDate, &book.CreatedAt, &book.UpdatedAt,
			&userEmail, &userCreatedAt,
		); err != nil {
			return nil, err
		}

		if userEmail != nil && userCreatedAt != nil {
			book.User = &domain.UserInfo{
				Email:     *userEmail,
				CreatedAt: *userCreatedAt,
			}
		}

		books = append(books, book)
	}
	return books, nil
}

func (r BookRepo) GetBookByID(id string) (*domain.Book, error) {
	var book domain.Book
	err := r.db.QueryRow(`
		SELECT id, title, author, genre, bookcase, user_email, return_date, created_at, updated_at
		FROM "book" WHERE id = $1
	`, id).Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Bookcase,
		&book.UserEmail, &book.ReturnDate, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r BookRepo) AssignBookToUser(bookID, userEmail string, returnDate time.Time) error {
	_, err := r.db.Exec(`
		UPDATE "book"
		SET user_email = $1, return_date = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, userEmail, returnDate, bookID)
	return err
}

func (r BookRepo) ReturnBook(bookID string) error {
	_, err := r.db.Exec(`
		UPDATE "book"
		SET user_email = NULL, return_date = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, bookID)
	return err
}

func (r BookRepo) DeleteBook(id string) error {
	_, err := r.db.Exec(`DELETE FROM "book" WHERE id = $1`, id)
	return err
}
