package types

type PasswordChangeRequest struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm"`
}

type PasswordChangeResetRequest struct {
	Email string `json:"email"`
}

type PasswordChangeResetConfirmRequest struct {
	Token              string `json:"token"`
	NewPassword        string `json:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm"`
}

type EmailConfirmationRequest struct {
	Token string `json:"token"`
}

type ResendEmailConfirmationRequest struct {
	Email string `json:"email"`
}

type UpdateUserRolesRequest struct {
	UserId string   `json:"user_id"`
	Roles  []string `json:"roles"`
}
