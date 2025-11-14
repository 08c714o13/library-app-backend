package services

import (
	"errors"
	"time"

	"library/internal/app/domain"
	"library/internal/app/repositories"

	"golang.org/x/crypto/bcrypt"
)

type LibrarianService struct {
	librarianRepo repositories.LibrarianRepoInterface
	userRepo      repositories.UserRepoInterface
	bookRepo      repositories.BookRepoInterface
}

func NewLibrarianService(librarianRepo repositories.LibrarianRepoInterface, userRepo repositories.UserRepoInterface, bookRepo repositories.BookRepoInterface) LibrarianService {
	return LibrarianService{
		librarianRepo: librarianRepo,
		userRepo:      userRepo,
		bookRepo:      bookRepo,
	}
}

func (s LibrarianService) Login(email string, password string) (bool, error) {
	hashPassword, err := s.librarianRepo.GetLibrarianByEmail(email)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return false, errors.New("invalid credentials")
	}

	return true, nil
}

func (s LibrarianService) CreateUser(email string) error {
	exists, err := s.userRepo.UserExists(email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	return s.userRepo.CreateUser(email)
}

func (s LibrarianService) DeleteUser(email string) error {
	exists, err := s.userRepo.UserExists(email)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	return s.userRepo.DeleteUser(email)
}

func (s LibrarianService) GetAllUsers() ([]domain.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s LibrarianService) GetAllUsersWithBooks() ([]domain.UserWithBooks, error) {
	return s.userRepo.GetAllUsersWithBooks()
}

func (s LibrarianService) CreateBook(book domain.CreateBookPayload) (string, error) {
	return s.bookRepo.CreateBook(book)
}

func (s LibrarianService) DeleteBook(id string) error {
	_, err := s.bookRepo.GetBookByID(id)
	if err != nil {
		return errors.New("book not found")
	}

	return s.bookRepo.DeleteBook(id)
}

func (s LibrarianService) GetAllBooks() ([]domain.Book, error) {
	return s.bookRepo.GetAllBooks()
}

func (s LibrarianService) GetAllBooksWithUsers() ([]domain.BookWithUser, error) {
	return s.bookRepo.GetAllBooksWithUsers()
}

func (s LibrarianService) AssignBookToUser(bookID string, userEmail string, returnDateStr string) error {
	book, err := s.bookRepo.GetBookByID(bookID)
	if err != nil {
		return errors.New("book not found")
	}

	if book.UserEmail != nil {
		return errors.New("book is already assigned to a user")
	}

	exists, err := s.userRepo.UserExists(userEmail)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	returnDate, err := time.Parse("02-01-2006", returnDateStr)
	if err != nil {
		return errors.New("invalid return date format, use DD-MM-YYYY")
	}
	if !returnDate.After(time.Now()) {
		return errors.New("")
	}

	return s.bookRepo.AssignBookToUser(bookID, userEmail, returnDate)
}

func (s LibrarianService) ReturnBook(bookID string) error {
	book, err := s.bookRepo.GetBookByID(bookID)
	if err != nil {
		return errors.New("book not found")
	}

	if book.UserEmail == nil {
		return errors.New("book is not assigned to any user")
	}

	return s.bookRepo.ReturnBook(bookID)
}
