package domain

import "time"

type Book struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Author     string     `json:"author"`
	Genre      string     `json:"genre"`
	Bookcase   string     `json:"bookcase"`
	UserEmail  *string    `json:"user_email,omitempty"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type BookWithUser struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Author     string     `json:"author"`
	Genre      string     `json:"genre"`
	Bookcase   string     `json:"bookcase"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	User       *UserInfo  `json:"user,omitempty"`
}

type CreateBookPayload struct {
	Title    string `json:"title" binding:"required"`
	Author   string `json:"author" binding:"required"`
	Genre    string `json:"genre" binding:"required"`
	Bookcase string `json:"bookcase" binding:"required"`
}

type AssignBookPayload struct {
	UserEmail  string `json:"user_email" binding:"required,email"`
	ReturnDate string `json:"return_date" binding:"required"`
}
