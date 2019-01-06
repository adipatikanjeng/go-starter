package requests

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
