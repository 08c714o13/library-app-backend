package domain

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Expire  string `json:"expire"`
	Message string `json:"message,omitempty"`
}
