package user

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Age      int8   `json:"age"`
	Gender   bool   `json:"gender"`
	Phone    string `json:"phone"`
}
