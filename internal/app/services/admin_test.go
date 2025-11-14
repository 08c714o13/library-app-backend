package services

import (
	"errors"
	"os"
	"testing"

	"library/internal/app/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockAdminRepo struct {
	mock.Mock
}

func (m *MockAdminRepo) CreateAdmin(email string, hashPassword string) error {
	args := m.Called(email, hashPassword)
	return args.Error(0)
}

func (m *MockAdminRepo) GetAdminByEmail(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

func (m *MockAdminRepo) AdminExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

type MockLibrarianRepoForAdmin struct {
	mock.Mock
}

func (m *MockLibrarianRepoForAdmin) LibrarianExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockLibrarianRepoForAdmin) CreateLibrarian(email, name, password string) error {
	args := m.Called(email, name, password)
	return args.Error(0)
}

func (m *MockLibrarianRepoForAdmin) DeleteLibrarian(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockLibrarianRepoForAdmin) GetAllLibrarians() ([]domain.Librarian, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Librarian), args.Error(1)
}

func (m *MockLibrarianRepoForAdmin) GetLibrarianByEmail(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

// Tests for Login
func TestAdminService_Login_Success(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	result, err := service.Login("admin@library.com", "admin123")

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestAdminService_Login_Success_WithEnvVars(t *testing.T) {
	os.Setenv("ADMIN_EMAIL", "test@admin.com")
	os.Setenv("ADMIN_PASSWORD", "testpass")
	defer os.Unsetenv("ADMIN_EMAIL")
	defer os.Unsetenv("ADMIN_PASSWORD")

	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	result, err := service.Login("test@admin.com", "testpass")

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestAdminService_Login_InvalidEmail(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	result, err := service.Login("wrong@email.com", "admin123")

	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestAdminService_Login_InvalidPassword(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	result, err := service.Login("admin@library.com", "wrongpass")

	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, "invalid credentials", err.Error())
}

// Tests for CreateLibrarian
func TestAdminService_CreateLibrarian_Success(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	librarianRepo.On("LibrarianExists", "test@lib.com").Return(false, nil)
	librarianRepo.On("CreateLibrarian", "test@lib.com", "Test Librarian", mock.AnythingOfType("string")).Return(nil)

	err := service.CreateLibrarian("test@lib.com", "Test Librarian", "password123")

	assert.NoError(t, err)
	librarianRepo.AssertExpectations(t)
}

func TestAdminService_CreateLibrarian_AlreadyExists(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	librarianRepo.On("LibrarianExists", "existing@lib.com").Return(true, nil)

	err := service.CreateLibrarian("existing@lib.com", "Test Librarian", "password123")

	assert.Error(t, err)
	assert.Equal(t, "librarian with this email already exists", err.Error())
	librarianRepo.AssertExpectations(t)
}

func TestAdminService_CreateLibrarian_CheckExistsError(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	librarianRepo.On("LibrarianExists", "test@lib.com").Return(false, errors.New("database error"))

	err := service.CreateLibrarian("test@lib.com", "Test Librarian", "password123")

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	librarianRepo.AssertExpectations(t)
}

func TestAdminService_CreateLibrarian_CreateError(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	librarianRepo.On("LibrarianExists", "test@lib.com").Return(false, nil)
	librarianRepo.On("CreateLibrarian", "test@lib.com", "Test Librarian", mock.AnythingOfType("string")).Return(errors.New("creation error"))

	err := service.CreateLibrarian("test@lib.com", "Test Librarian", "password123")

	assert.Error(t, err)
	assert.Equal(t, "creation error", err.Error())
	librarianRepo.AssertExpectations(t)
}

// Tests for DeleteLibrarian
func TestAdminService_DeleteLibrarian_Success(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	librarianRepo.On("LibrarianExists", "test@lib.com").Return(true, nil)
	librarianRepo.On("DeleteLibrarian", "test@lib.com").Return(nil)

	err := service.DeleteLibrarian("test@lib.com")

	assert.NoError(t, err)
	librarianRepo.AssertExpectations(t)
}

func TestAdminService_DeleteLibrarian_NotFound(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	librarianRepo.On("LibrarianExists", "nonexistent@lib.com").Return(false, nil)

	err := service.DeleteLibrarian("nonexistent@lib.com")

	assert.Error(t, err)
	assert.Equal(t, "librarian not found", err.Error())
	librarianRepo.AssertExpectations(t)
}

func TestAdminService_DeleteLibrarian_CheckExistsError(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	librarianRepo.On("LibrarianExists", "test@lib.com").Return(false, errors.New("database error"))

	err := service.DeleteLibrarian("test@lib.com")

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	librarianRepo.AssertExpectations(t)
}

func TestAdminService_DeleteLibrarian_DeleteError(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	librarianRepo.On("LibrarianExists", "test@lib.com").Return(true, nil)
	librarianRepo.On("DeleteLibrarian", "test@lib.com").Return(errors.New("deletion error"))

	err := service.DeleteLibrarian("test@lib.com")

	assert.Error(t, err)
	assert.Equal(t, "deletion error", err.Error())
	librarianRepo.AssertExpectations(t)
}

// Tests for GetAllLibrarians
func TestAdminService_GetAllLibrarians_Success(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	mockLibrarians := []domain.Librarian{
		{Email: "lib1@test.com", Name: "Librarian 1"},
		{Email: "lib2@test.com", Name: "Librarian 2"},
	}
	librarianRepo.On("GetAllLibrarians").Return(mockLibrarians, nil)

	result, err := service.GetAllLibrarians()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	librarianRepo.AssertExpectations(t)
}

func TestAdminService_GetAllLibrarians_Empty(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	mockLibrarians := []domain.Librarian{}
	librarianRepo.On("GetAllLibrarians").Return(mockLibrarians, nil)

	result, err := service.GetAllLibrarians()

	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))
	librarianRepo.AssertExpectations(t)
}

func TestAdminService_GetAllLibrarians_Error(t *testing.T) {
	adminRepo := new(MockAdminRepo)
	librarianRepo := new(MockLibrarianRepoForAdmin)
	service := NewAdminService(adminRepo, librarianRepo)

	var nilLibrarians []domain.Librarian
	librarianRepo.On("GetAllLibrarians").Return(nilLibrarians, errors.New("database error"))

	result, err := service.GetAllLibrarians()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	librarianRepo.AssertExpectations(t)
}
