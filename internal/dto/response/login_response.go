package response

type LoginUserResponse struct {
	Token    string `json:"token"`
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
