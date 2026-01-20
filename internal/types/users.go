package types

type PasswordChangeRequest struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm"`
}

type UpdateUserRolesRequest struct {
	UserId string   `json:"user_id"`
	Roles  []string `json:"roles"`
}
