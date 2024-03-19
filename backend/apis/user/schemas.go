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
	Age      int8   `json:"age" validate:"omitempty"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}
