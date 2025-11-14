package controllers

import (
	"net/http"

	"library/internal/app/domain"
	"library/internal/app/services"

	"github.com/gin-gonic/gin"
)

type LibrarianController struct {
	librarianService services.LibrarianService
}

func NewLibrarianController(librarianService services.LibrarianService) LibrarianController {
	return LibrarianController{librarianService: librarianService}
}

// CreateUser godoc
//
//	@Summary		Создать пользователя
//	@Description	Создает нового пользователя библиотеки
//	@Tags			Librarian - Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		domain.CreateUserPayload	true	"Email пользователя"
//	@Success		201		{object}	map[string]interface{}		"User created successfully"
//	@Failure		400		{object}	map[string]interface{}		"Invalid request"
//	@Security		Bearer
//	@Router			/librarian/users [post]
func (lc LibrarianController) CreateUser(c *gin.Context) {
	var payload domain.CreateUserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if err := lc.librarianService.CreateUser(payload.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "email": payload.Email})
}

// DeleteUser godoc
//
//	@Summary		Удалить пользователя
//	@Description	Удаляет пользователя по email
//	@Tags			Librarian - Users
//	@Produce		json
//	@Param			email	path		string					true	"Email пользователя"
//	@Success		200		{object}	map[string]interface{}	"User deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Error message"
//	@Security		Bearer
//	@Router			/librarian/users/{email} [delete]
func (lc LibrarianController) DeleteUser(c *gin.Context) {
	email := c.Param("email")

	if err := lc.librarianService.DeleteUser(email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetAllUsers godoc
//
//	@Summary		Получить всех пользователей
//	@Description	Возвращает список всех пользователей. Используйте ?with=books для получения пользователей с их книгами
//	@Tags			Librarian - Users
//	@Produce		json
//	@Param			with	query		string					false	"Включить дополнительные данные (books)"
//	@Success		200		{object}	map[string]interface{}	"List of users"
//	@Failure		500		{object}	map[string]interface{}	"Error message"
//	@Security		Bearer
//	@Router			/librarian/users [get]
func (lc LibrarianController) GetAllUsers(c *gin.Context) {
	// Check if we should include books
	withBooks := c.Query("with") == "books"

	if withBooks {
		users, err := lc.librarianService.GetAllUsersWithBooks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
		return
	}

	users, err := lc.librarianService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// CreateBook godoc
//
//	@Summary		Создать книгу
//	@Description	Добавляет новую книгу в библиотеку
//	@Tags			Librarian - Books
//	@Accept			json
//	@Produce		json
//	@Param			book	body		domain.CreateBookPayload	true	"Данные книги"
//	@Success		201		{object}	map[string]interface{}		"Book created successfully"
//	@Failure		400		{object}	map[string]interface{}		"Invalid request"
//	@Security		Bearer
//	@Router			/librarian/books [post]
func (lc LibrarianController) CreateBook(c *gin.Context) {
	var payload domain.CreateBookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	id, err := lc.librarianService.CreateBook(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book created successfully", "id": id})
}

// DeleteBook godoc
//
//	@Summary		Удалить книгу
//	@Description	Удаляет книгу по ID
//	@Tags			Librarian - Books
//	@Produce		json
//	@Param			id	path		string					true	"ID книги"
//	@Success		200	{object}	map[string]interface{}	"Book deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Error message"
//	@Security		Bearer
//	@Router			/librarian/books/{id} [delete]
func (lc LibrarianController) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	if err := lc.librarianService.DeleteBook(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// GetAllBooks godoc
//
//	@Summary		Получить все книги
//	@Description	Возвращает список всех книг. Используйте ?with=users для получения книг с информацией о пользователях
//	@Tags			Librarian - Books
//	@Produce		json
//	@Param			with	query		string					false	"Включить дополнительные данные (users)"
//	@Success		200		{object}	map[string]interface{}	"List of books"
//	@Failure		500		{object}	map[string]interface{}	"Error message"
//	@Security		Bearer
//	@Router			/librarian/books [get]
func (lc LibrarianController) GetAllBooks(c *gin.Context) {
	withUsers := c.Query("with") == "users"

	if withUsers {
		books, err := lc.librarianService.GetAllBooksWithUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"books": books})
		return
	}

	books, err := lc.librarianService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}

// AssignBook godoc
//
//	@Summary		Выдать книгу пользователю
//	@Description	Назначает книгу пользователю с датой возврата
//	@Tags			Librarian - Books
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"ID книги"
//	@Param			assign	body		domain.AssignBookPayload	true	"Данные назначения"
//	@Success		200		{object}	map[string]interface{}		"Book assigned successfully"
//	@Failure		400		{object}	map[string]interface{}		"Error message"
//	@Security		Bearer
//	@Router			/librarian/books/{id}/assign [post]
func (lc LibrarianController) AssignBook(c *gin.Context) {
	bookID := c.Param("id")
	var payload domain.AssignBookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if err := lc.librarianService.AssignBookToUser(bookID, payload.UserEmail, payload.ReturnDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book assigned successfully"})
}

// ReturnBook godoc
//
//	@Summary		Вернуть книгу
//	@Description	Отмечает книгу как возвращенную
//	@Tags			Librarian - Books
//	@Produce		json
//	@Param			id	path		string					true	"ID книги"
//	@Success		200	{object}	map[string]interface{}	"Book returned successfully"
//	@Failure		400	{object}	map[string]interface{}	"Error message"
//	@Security		Bearer
//	@Router			/librarian/books/{id}/return [post]
func (lc LibrarianController) ReturnBook(c *gin.Context) {
	bookID := c.Param("id")

	if err := lc.librarianService.ReturnBook(bookID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}
