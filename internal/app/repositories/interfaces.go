package repositories

import (
	"time"

	"library/internal/app/domain"
)

type AdminRepoInterface interface {
	CreateAdmin(email string, hashPassword string) error
	GetAdminByEmail(email string) (string, error)
	AdminExists(email string) (bool, error)
}

type LibrarianRepoInterface interface {
	CreateLibrarian(email, name, hashPassword string) error
	GetLibrarianByEmail(email string) (string, error)
	GetAllLibrarians() ([]domain.Librarian, error)
	DeleteLibrarian(email string) error
	LibrarianExists(email string) (bool, error)
}

type UserRepoInterface interface {
	CreateUser(email string) error
	GetAllUsers() ([]domain.User, error)
	GetAllUsersWithBooks() ([]domain.UserWithBooks, error)
	GetUserByEmail(email string) (*domain.User, error)
	DeleteUser(email string) error
	UserExists(email string) (bool, error)
}

type BookRepoInterface interface {
	CreateBook(book domain.CreateBookPayload) (string, error)
	GetAllBooks() ([]domain.Book, error)
	GetAllBooksWithUsers() ([]domain.BookWithUser, error)
	GetBookByID(id string) (*domain.Book, error)
	AssignBookToUser(bookID, userEmail string, returnDate time.Time) error
	ReturnBook(bookID string) error
	DeleteBook(id string) error
}
