package response

type RegisterResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}
