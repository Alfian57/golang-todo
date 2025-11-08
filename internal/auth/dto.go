package auth

// General UserResponse
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// Login
type LoginRequest struct {
	Username string `json:"username" validate:"required,max=255,min=1"`
	Password string `json:"password" validate:"required,max=100,min=1"`
}
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// Register
type RegisterRequest struct {
	Username             string `json:"username" validate:"required,max=255,min=1"`
	Password             string `json:"password" validate:"required,max=100,min=1"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,max=100,min=1,eqfield=Password"`
}
type RegisterResponse struct {
	User UserResponse `json:"user"`
}

// Logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
}
