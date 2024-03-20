package user

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Age      int8   `json:"age"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
}

// difference between pointer and value withomitempty
type UpdateUserRequest struct {
	Age    int8   `json:"age" validate:"omitempty"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"omitempty"`
	Email       string `json:"email" validate:"omitempty"`
	NewPassword string `json:"new_password"`
}

type TokenResponse struct {
	Access  string `json:"access,omitempty"`
	Message string `json:"message,omitempty"`
}
