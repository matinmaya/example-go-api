package usermodel

type AuthChangePassword struct {
	OldPassword     string `json:"old_password" validate:"required,min=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
}

type AuthCredentials struct {
	Username string `json:"username" validate:"required,min=6"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthLoginResource struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
