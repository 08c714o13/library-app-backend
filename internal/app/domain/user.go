package domain

import "time"

type User struct {
	Email     string    `json:"email"`
	Role      string    `json:"role,omitempty"`
	IsDeleted bool      `json:"is_deleted,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UserInfo struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserWithBooks struct {
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Books     []BookInfo `json:"books,omitempty"`
}

type BookInfo struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Author     string     `json:"author"`
	Genre      string     `json:"genre"`
	Bookcase   string     `json:"bookcase"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
}
