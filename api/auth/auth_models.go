package auth

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDto struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
