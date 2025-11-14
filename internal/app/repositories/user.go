package repositories

import (
	"database/sql"

	"library/internal/app/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{db: db}
}

func (r UserRepo) CreateUser(email string) error {
	_, err := r.db.Exec(`INSERT INTO "user" (email) VALUES ($1)`, email)
	return err
}

func (r UserRepo) GetAllUsers() ([]domain.User, error) {
	rows, err := r.db.Query(`SELECT email, is_deleted, created_at, updated_at FROM "user" WHERE is_deleted = false ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.Email, &user.IsDeleted, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r UserRepo) GetAllUsersWithBooks() ([]domain.UserWithBooks, error) {
	userRows, err := r.db.Query(`
		SELECT email, created_at, updated_at
		FROM "user"
		WHERE is_deleted = false
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer func(userRows *sql.Rows) {
		_ = userRows.Close()
	}(userRows)

	var users []domain.UserWithBooks
	userMap := make(map[string]*domain.UserWithBooks)

	for userRows.Next() {
		var user domain.UserWithBooks
		if err := userRows.Scan(&user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		user.Books = []domain.BookInfo{}
		users = append(users, user)
		userMap[user.Email] = &users[len(users)-1]
	}

	bookRows, err := r.db.Query(`
		SELECT
			b.id, b.title, b.author, b.genre, b.bookcase, b.return_date, b.user_email
		FROM "book" b
		WHERE b.user_email IS NOT NULL
		ORDER BY b.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer func(bookRows *sql.Rows) {
		_ = bookRows.Close()
	}(bookRows)

	for bookRows.Next() {
		var book domain.BookInfo
		var userEmail string

		if err := bookRows.Scan(
			&book.ID, &book.Title, &book.Author, &book.Genre, &book.Bookcase,
			&book.ReturnDate, &userEmail,
		); err != nil {
			return nil, err
		}

		if user, exists := userMap[userEmail]; exists {
			user.Books = append(user.Books, book)
		}
	}

	return users, nil
}

func (r UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`SELECT email, is_deleted, created_at, updated_at FROM "user" WHERE email = $1`, email).
		Scan(&user.Email, &user.IsDeleted, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r UserRepo) DeleteUser(email string) error {
	_, err := r.db.Exec(`UPDATE "user" SET is_deleted = true, updated_at = CURRENT_TIMESTAMP WHERE email = $1`, email)
	return err
}

func (r UserRepo) UserExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1 AND is_deleted = false)`, email).Scan(&exists)
	return exists, err
}
