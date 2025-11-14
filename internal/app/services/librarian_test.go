package services

import (
	"errors"
	"testing"
	"time"

	"library/internal/app/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock repositories
type MockLibrarianRepo struct {
	mock.Mock
}

func (m *MockLibrarianRepo) CreateLibrarian(email, name, hashPassword string) error {
	args := m.Called(email, name, hashPassword)
	return args.Error(0)
}

func (m *MockLibrarianRepo) GetLibrarianByEmail(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

func (m *MockLibrarianRepo) GetAllLibrarians() ([]domain.Librarian, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Librarian), args.Error(1)
}

func (m *MockLibrarianRepo) DeleteLibrarian(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockLibrarianRepo) LibrarianExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

type MockUserRepo struct {
	mock.Mock
}

type MockBookRepo struct {
	mock.Mock
}

func (m *MockUserRepo) UserExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) CreateUser(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserRepo) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepo) GetAllUsersWithBooks() ([]domain.UserWithBooks, error) {
	args := m.Called()
	return args.Get(0).([]domain.UserWithBooks), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockBookRepo) CreateBook(book domain.CreateBookPayload) (string, error) {
	args := m.Called(book)
	return args.String(0), args.Error(1)
}

func (m *MockBookRepo) DeleteBook(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookRepo) GetAllBooks() ([]domain.Book, error) {
	args := m.Called()
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (m *MockBookRepo) GetAllBooksWithUsers() ([]domain.BookWithUser, error) {
	args := m.Called()
	return args.Get(0).([]domain.BookWithUser), args.Error(1)
}

func (m *MockBookRepo) GetBookByID(id string) (*domain.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookRepo) AssignBookToUser(bookID string, userEmail string, returnDate time.Time) error {
	args := m.Called(bookID, userEmail, returnDate)
	return args.Error(0)
}

func (m *MockBookRepo) ReturnBook(bookID string) error {
	args := m.Called(bookID)
	return args.Error(0)
}

// Tests for Login
func TestLibrarianService_Login_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	librarianRepo.On("GetLibrarianByEmail", "lib@test.com").Return(string(hashedPassword), nil)

	result, err := service.Login("lib@test.com", password)

	assert.NoError(t, err)
	assert.True(t, result)
	librarianRepo.AssertExpectations(t)
}

func TestLibrarianService_Login_InvalidPassword(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	librarianRepo.On("GetLibrarianByEmail", "lib@test.com").Return(string(hashedPassword), nil)

	result, err := service.Login("lib@test.com", "wrongpassword")

	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, "invalid credentials", err.Error())
	librarianRepo.AssertExpectations(t)
}

func TestLibrarianService_Login_LibrarianNotFound(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	librarianRepo.On("GetLibrarianByEmail", "nonexistent@test.com").Return("", errors.New("librarian not found"))

	result, err := service.Login("nonexistent@test.com", "password123")

	assert.Error(t, err)
	assert.False(t, result)
	librarianRepo.AssertExpectations(t)
}

// Tests for CreateUser
func TestLibrarianService_CreateUser_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	userRepo.On("UserExists", "user@test.com").Return(false, nil)
	userRepo.On("CreateUser", "user@test.com").Return(nil)

	err := service.CreateUser("user@test.com")

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestLibrarianService_CreateUser_AlreadyExists(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	userRepo.On("UserExists", "existing@test.com").Return(true, nil)

	err := service.CreateUser("existing@test.com")

	assert.Error(t, err)
	assert.Equal(t, "user with this email already exists", err.Error())
	userRepo.AssertExpectations(t)
}

func TestLibrarianService_CreateUser_CheckExistsError(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	userRepo.On("UserExists", "user@test.com").Return(false, errors.New("database error"))

	err := service.CreateUser("user@test.com")

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	userRepo.AssertExpectations(t)
}

// Tests for DeleteUser
func TestLibrarianService_DeleteUser_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	userRepo.On("UserExists", "user@test.com").Return(true, nil)
	userRepo.On("DeleteUser", "user@test.com").Return(nil)

	err := service.DeleteUser("user@test.com")

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestLibrarianService_DeleteUser_NotFound(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	userRepo.On("UserExists", "nonexistent@test.com").Return(false, nil)

	err := service.DeleteUser("nonexistent@test.com")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}

// Tests for GetAllUsers
func TestLibrarianService_GetAllUsers_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockUsers := []domain.User{
		{Email: "user1@test.com"},
		{Email: "user2@test.com"},
	}
	userRepo.On("GetAllUsers").Return(mockUsers, nil)

	result, err := service.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	userRepo.AssertExpectations(t)
}

func TestLibrarianService_GetAllUsers_Error(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	var nilUsers []domain.User
	userRepo.On("GetAllUsers").Return(nilUsers, errors.New("database error"))

	result, err := service.GetAllUsers()

	assert.Error(t, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
}

// Tests for GetAllUsersWithBooks
func TestLibrarianService_GetAllUsersWithBooks_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockUsersWithBooks := []domain.UserWithBooks{
		{Email: "user1@test.com", Books: []domain.BookInfo{{ID: "1", Title: "Book 1"}}},
	}
	userRepo.On("GetAllUsersWithBooks").Return(mockUsersWithBooks, nil)

	result, err := service.GetAllUsersWithBooks()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	userRepo.AssertExpectations(t)
}

// Tests for CreateBook
func TestLibrarianService_CreateBook_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	bookPayload := domain.CreateBookPayload{
		Title:    "Test Book",
		Author:   "Test Author",
		Genre:    "Fiction",
		Bookcase: "A1",
	}

	bookRepo.On("CreateBook", bookPayload).Return("book-id-123", nil)

	bookID, err := service.CreateBook(bookPayload)

	assert.NoError(t, err)
	assert.Equal(t, "book-id-123", bookID)
	bookRepo.AssertExpectations(t)
}

func TestLibrarianService_CreateBook_Error(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	bookPayload := domain.CreateBookPayload{
		Title:    "Test Book",
		Author:   "Test Author",
		Genre:    "Fiction",
		Bookcase: "A1",
	}

	bookRepo.On("CreateBook", bookPayload).Return("", errors.New("creation error"))

	bookID, err := service.CreateBook(bookPayload)

	assert.Error(t, err)
	assert.Equal(t, "", bookID)
	bookRepo.AssertExpectations(t)
}

// Tests for DeleteBook
func TestLibrarianService_DeleteBook_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockBook := &domain.Book{ID: "book-123", Title: "Test Book"}
	bookRepo.On("GetBookByID", "book-123").Return(mockBook, nil)
	bookRepo.On("DeleteBook", "book-123").Return(nil)

	err := service.DeleteBook("book-123")

	assert.NoError(t, err)
	bookRepo.AssertExpectations(t)
}

func TestLibrarianService_DeleteBook_NotFound(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	var nilBook *domain.Book
	bookRepo.On("GetBookByID", "nonexistent").Return(nilBook, errors.New("not found"))

	err := service.DeleteBook("nonexistent")

	assert.Error(t, err)
	assert.Equal(t, "book not found", err.Error())
	bookRepo.AssertExpectations(t)
}

// Tests for GetAllBooks
func TestLibrarianService_GetAllBooks_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockBooks := []domain.Book{
		{ID: "1", Title: "Book 1"},
		{ID: "2", Title: "Book 2"},
	}
	bookRepo.On("GetAllBooks").Return(mockBooks, nil)

	result, err := service.GetAllBooks()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	bookRepo.AssertExpectations(t)
}

// Tests for GetAllBooksWithUsers
func TestLibrarianService_GetAllBooksWithUsers_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockBooksWithUsers := []domain.BookWithUser{
		{ID: "1", Title: "Book 1", User: &domain.UserInfo{Email: "user@test.com"}},
	}
	bookRepo.On("GetAllBooksWithUsers").Return(mockBooksWithUsers, nil)

	result, err := service.GetAllBooksWithUsers()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	bookRepo.AssertExpectations(t)
}

// Tests for AssignBookToUser
func TestLibrarianService_AssignBookToUser_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockBook := &domain.Book{ID: "book-123", Title: "Test Book", UserEmail: nil}
	futureDate := time.Now().Add(7 * 24 * time.Hour)
	returnDateStr := futureDate.Format("02-01-2006")

	bookRepo.On("GetBookByID", "book-123").Return(mockBook, nil)
	userRepo.On("UserExists", "user@test.com").Return(true, nil)
	bookRepo.On("AssignBookToUser", "book-123", "user@test.com", mock.AnythingOfType("time.Time")).Return(nil)

	err := service.AssignBookToUser("book-123", "user@test.com", returnDateStr)

	assert.NoError(t, err)
	bookRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestLibrarianService_AssignBookToUser_BookNotFound(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	var nilBook *domain.Book
	bookRepo.On("GetBookByID", "nonexistent").Return(nilBook, errors.New("not found"))

	err := service.AssignBookToUser("nonexistent", "user@test.com", "01-01-2025")

	assert.Error(t, err)
	assert.Equal(t, "book not found", err.Error())
	bookRepo.AssertExpectations(t)
}

func TestLibrarianService_AssignBookToUser_BookAlreadyAssigned(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	assignedEmail := "other@test.com"
	mockBook := &domain.Book{ID: "book-123", Title: "Test Book", UserEmail: &assignedEmail}

	bookRepo.On("GetBookByID", "book-123").Return(mockBook, nil)

	err := service.AssignBookToUser("book-123", "user@test.com", "01-01-2025")

	assert.Error(t, err)
	assert.Equal(t, "book is already assigned to a user", err.Error())
	bookRepo.AssertExpectations(t)
}

func TestLibrarianService_AssignBookToUser_UserNotFound(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockBook := &domain.Book{ID: "book-123", Title: "Test Book", UserEmail: nil}

	bookRepo.On("GetBookByID", "book-123").Return(mockBook, nil)
	userRepo.On("UserExists", "nonexistent@test.com").Return(false, nil)

	err := service.AssignBookToUser("book-123", "nonexistent@test.com", "01-01-2025")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	bookRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestLibrarianService_AssignBookToUser_InvalidDateFormat(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockBook := &domain.Book{ID: "book-123", Title: "Test Book", UserEmail: nil}

	bookRepo.On("GetBookByID", "book-123").Return(mockBook, nil)
	userRepo.On("UserExists", "user@test.com").Return(true, nil)

	err := service.AssignBookToUser("book-123", "user@test.com", "invalid-date")

	assert.Error(t, err)
	assert.Equal(t, "invalid return date format, use DD-MM-YYYY", err.Error())
	bookRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestLibrarianService_AssignBookToUser_PastDate(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockBook := &domain.Book{ID: "book-123", Title: "Test Book", UserEmail: nil}
	pastDate := time.Now().Add(-7 * 24 * time.Hour)
	returnDateStr := pastDate.Format("02-01-2006")

	bookRepo.On("GetBookByID", "book-123").Return(mockBook, nil)
	userRepo.On("UserExists", "user@test.com").Return(true, nil)

	err := service.AssignBookToUser("book-123", "user@test.com", returnDateStr)

	assert.Error(t, err)
	bookRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

// Tests for ReturnBook
func TestLibrarianService_ReturnBook_Success(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	assignedEmail := "user@test.com"
	mockBook := &domain.Book{ID: "book-123", Title: "Test Book", UserEmail: &assignedEmail}

	bookRepo.On("GetBookByID", "book-123").Return(mockBook, nil)
	bookRepo.On("ReturnBook", "book-123").Return(nil)

	err := service.ReturnBook("book-123")

	assert.NoError(t, err)
	bookRepo.AssertExpectations(t)
}

func TestLibrarianService_ReturnBook_BookNotFound(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	var nilBook *domain.Book
	bookRepo.On("GetBookByID", "nonexistent").Return(nilBook, errors.New("not found"))

	err := service.ReturnBook("nonexistent")

	assert.Error(t, err)
	assert.Equal(t, "book not found", err.Error())
	bookRepo.AssertExpectations(t)
}

func TestLibrarianService_ReturnBook_BookNotAssigned(t *testing.T) {
	librarianRepo := new(MockLibrarianRepo)
	userRepo := new(MockUserRepo)
	bookRepo := new(MockBookRepo)
	service := NewLibrarianService(librarianRepo, userRepo, bookRepo)

	mockBook := &domain.Book{ID: "book-123", Title: "Test Book", UserEmail: nil}

	bookRepo.On("GetBookByID", "book-123").Return(mockBook, nil)

	err := service.ReturnBook("book-123")

	assert.Error(t, err)
	assert.Equal(t, "book is not assigned to any user", err.Error())
	bookRepo.AssertExpectations(t)
}
