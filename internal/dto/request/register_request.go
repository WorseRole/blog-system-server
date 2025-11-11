package request

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" bingding:"required,min=6, max=16"`
	Email    string `json:"email" binding:"required,email"`
}
