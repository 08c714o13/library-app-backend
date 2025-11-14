package controllers

import (
	"net/http"

	"library/internal/app/domain"
	"library/internal/app/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService services.AdminService
}

func NewAdminController(adminService services.AdminService) AdminController {
	return AdminController{adminService: adminService}
}

// CreateLibrarian godoc
//
//	@Summary		Создать библиотекаря
//	@Description	Создает нового библиотекаря (только для администратора)
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			librarian	body		domain.CreateLibrarianPayload	true	"Данные библиотекаря"
//	@Success		201			{object}	map[string]interface{}			"Librarian created successfully"
//	@Failure		400			{object}	map[string]interface{}			"Invalid request"
//	@Security		Bearer
//	@Router			/admin/librarians [post]
func (a AdminController) CreateLibrarian(c *gin.Context) {
	var payload domain.CreateLibrarianPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if err := a.adminService.CreateLibrarian(payload.Email, payload.Name, payload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Librarian created successfully", "email": payload.Email})
}

// DeleteLibrarian godoc
//
//	@Summary		Удалить библиотекаря
//	@Description	Удаляет библиотекаря по email (только для администратора)
//	@Tags			Admin
//	@Produce		json
//	@Param			email	path		string					true	"Email библиотекаря"
//	@Success		200		{object}	map[string]interface{}	"Librarian deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Error message"
//	@Security		Bearer
//	@Router			/admin/librarians/{email} [delete]
func (a AdminController) DeleteLibrarian(c *gin.Context) {
	email := c.Param("email")

	if err := a.adminService.DeleteLibrarian(email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Librarian deleted successfully"})
}

// GetAllLibrarians godoc
//
//	@Summary		Получить всех библиотекарей
//	@Description	Возвращает список всех библиотекарей (только для администратора)
//	@Tags			Admin
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"List of librarians"
//	@Failure		500	{object}	map[string]interface{}	"Error message"
//	@Security		Bearer
//	@Router			/admin/librarians [get]
func (a AdminController) GetAllLibrarians(c *gin.Context) {
	librarians, err := a.adminService.GetAllLibrarians()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"librarians": librarians})
}
