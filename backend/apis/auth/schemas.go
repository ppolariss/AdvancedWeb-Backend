package auth

//type EmailModel struct {
//	// email in email blacklist
//	Email string `json:"email" query:"email" validate:"isValidEmail"`
//}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" minLength:"3" validate:"required,min=3"`
}

type TokenResponse struct {
	Access  string `json:"access,omitempty"`
	Message string `json:"message,omitempty"`
	UserID  int    `json:"user_id,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" minLength:"3" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Age      int8   `json:"age" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}
