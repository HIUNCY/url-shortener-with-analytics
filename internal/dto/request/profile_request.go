package request

// UpdateProfileRequest adalah DTO untuk request pembaruan profil.
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// ChangePasswordRequest adalah DTO untuk request perubahan password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}
