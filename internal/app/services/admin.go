package services

import (
	"errors"
	"os"

	"library/internal/app/domain"
	"library/internal/app/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	adminRepo     repositories.AdminRepoInterface
	librarianRepo repositories.LibrarianRepoInterface
}

func NewAdminService(adminRepo repositories.AdminRepoInterface, librarianRepo repositories.LibrarianRepoInterface) AdminService {
	return AdminService{adminRepo: adminRepo, librarianRepo: librarianRepo}
}

func (a AdminService) Login(email string, password string) (bool, error) {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@library.com"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	if email != adminEmail || password != adminPassword {
		return false, errors.New("invalid credentials")
	}

	return true, nil
}

func (a AdminService) CreateLibrarian(email, name, password string) error {
	exists, err := a.librarianRepo.LibrarianExists(email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("librarian with this email already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return a.librarianRepo.CreateLibrarian(email, name, string(hashPassword))
}

func (a AdminService) DeleteLibrarian(email string) error {
	exists, err := a.librarianRepo.LibrarianExists(email)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("librarian not found")
	}

	return a.librarianRepo.DeleteLibrarian(email)
}

func (a AdminService) GetAllLibrarians() ([]domain.Librarian, error) {
	return a.librarianRepo.GetAllLibrarians()
}
