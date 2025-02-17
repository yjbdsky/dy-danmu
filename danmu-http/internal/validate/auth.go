package validate

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty,min=6,max=50"`
	Role     string `json:"role" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateSelfRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6,max=50"`
}
