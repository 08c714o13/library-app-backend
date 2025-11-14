package main

import (
	"database/sql"
	"log"
	"os"

	"library/cmd/http"
	"library/internal/app/controllers"
	"library/internal/app/repositories"
	"library/internal/app/services"
	"library/internal/pkg/database"
	"library/internal/pkg/middlewares"

	"github.com/joho/godotenv"
)

//	@title			Library Management System API
//	@version		1.0
//	@description	API для управления библиотекой с ролями администратора и библиотекаря
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.email	admin@library.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	loadEnv()

	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "library"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	migrationsDir := getEnv("MIGRATIONS_DIR", "./migration")
	if err := database.RunMigrations(db, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	adminRepo := repositories.NewAdminRepo(db)
	librarianRepo := repositories.NewLibrarianRepo(db)
	userRepo := repositories.NewUserRepo(db)
	bookRepo := repositories.NewBookRepo(db)

	adminService := services.NewAdminService(adminRepo, librarianRepo)
	librarianService := services.NewLibrarianService(librarianRepo, userRepo, bookRepo)

	adminController := controllers.NewAdminController(adminService)
	librarianController := controllers.NewLibrarianController(librarianService)

	authMiddleware, err := middlewares.InitAuthJWT(adminService, librarianService)
	if err != nil {
		log.Fatalf("Failed to initialize JWT middleware: %v", err)
	}

	router := http.CreateRoutes(adminController, librarianController, authMiddleware)

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
		return
	}
	log.Println(".env file loaded successfully")
}
