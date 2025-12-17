package repositories

import (
	"database/sql"
	"time"

	"library/internal/app/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{db: db}
}

func (r UserRepo) CreateUser(email string) error {
	_, err := r.db.Exec(`
    INSERT INTO "user" (email, is_deleted, created_at, updated_at) 
    VALUES ($1, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    ON CONFLICT (email) 
    DO UPDATE SET 
        is_deleted = false,
        updated_at = CURRENT_TIMESTAMP
`, email)
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

func (r UserRepo) GetAllUsersWithBooks() ([]domain.UserWithBooks, error) { // TODO
	rows, err := r.db.Query(`
    SELECT 
        u.email, u.created_at, u.updated_at,
        b.id, b.title, b.author, b.genre, b.bookcase, b.return_date
    FROM "user" u
    LEFT JOIN "book" b ON u.email = b.user_email AND b.user_email IS NOT NULL
    ORDER BY u.created_at DESC, b.created_at DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.UserWithBooks
	userIndexMap := make(map[string]int)

	for rows.Next() {
		var email string
		var createdAt, updatedAt time.Time
		var bookID sql.NullString
		var bookTitle, bookAuthor, bookGenre, bookBookcase sql.NullString
		var returnDate sql.NullTime

		if err := rows.Scan(
			&email, &createdAt, &updatedAt,
			&bookID, &bookTitle, &bookAuthor, &bookGenre, &bookBookcase, &returnDate,
		); err != nil {
			return nil, err
		}

		// Создаем пользователя, если его еще нет в мапе
		idx, exists := userIndexMap[email]
		if !exists {
			user := domain.UserWithBooks{
				Email:     email,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				Books:     []domain.BookInfo{},
			}
			users = append(users, user)
			idx = len(users) - 1
			userIndexMap[email] = idx
		}

		// Добавляем книгу, если она есть
		if bookID.Valid {
			book := domain.BookInfo{
				ID:         bookID.String,
				Title:      bookTitle.String,
				Author:     bookAuthor.String,
				Genre:      bookGenre.String,
				Bookcase:   bookBookcase.String,
				ReturnDate: &returnDate.Time,
			}
			users[idx].Books = append(users[idx].Books, book)
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
