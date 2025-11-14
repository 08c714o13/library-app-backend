package http

import (
	"library/internal/app/controllers"
	"library/internal/pkg/middlewares"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateRoutes(
	adminController controllers.AdminController,
	librarianController controllers.LibrarianController,
	authMiddleware *jwt.GinJWTMiddleware,
) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api")

	public.POST("/admin/login", authMiddleware.LoginHandler)
	public.POST("/librarian/login", authMiddleware.LoginHandler)

	public.GET("/refresh_token", authMiddleware.RefreshHandler)

	admin := r.Group("/api/admin")
	admin.Use(authMiddleware.MiddlewareFunc())
	admin.Use(middlewares.RoleMiddleware("admin"))

	admin.POST("/librarians", adminController.CreateLibrarian)
	admin.DELETE("/librarians/:email", adminController.DeleteLibrarian)
	admin.GET("/librarians", adminController.GetAllLibrarians)

	librarian := r.Group("/api/librarian")
	librarian.Use(authMiddleware.MiddlewareFunc())
	librarian.Use(middlewares.RoleMiddleware("librarian"))

	librarian.POST("/users", librarianController.CreateUser)
	librarian.DELETE("/users/:email", librarianController.DeleteUser)
	librarian.GET("/users", librarianController.GetAllUsers)

	librarian.POST("/books", librarianController.CreateBook)
	librarian.DELETE("/books/:id", librarianController.DeleteBook)
	librarian.GET("/books", librarianController.GetAllBooks)

	librarian.POST("/books/:id/assign", librarianController.AssignBook)
	librarian.POST("/books/:id/return", librarianController.ReturnBook)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
