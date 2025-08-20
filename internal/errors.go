package internal

import "fmt"

const (
	ErrAccessNotPermitted        = "access not permitted"
	ErrUserNotFound              = "user not found"
	ErrPermitNotFound            = "permit not found"
	ErrRoleNotFound              = "role not found"
	ErrPermissionNotFound        = "permission not found"
	ErrGrantingPermit            = "error granting permit"
	ErrRevokingPermit            = "error revoking permit"
	ErrUnmarshalingPermits       = "error unmarshaling permits"
	ErrMarshallingPermits        = "error marshalling permits"
	ErrGettingRolePermissions    = "error getting role permissions"
	ErrRoleOrPermissionNotFound  = "role or permission not found"
	ErrAddingRolePermission      = "error adding role permission"
	ErrDeletingRolePermission    = "error deleting role permission"
	ErrGettingRoles              = "error getting roles"
	ErrAddingRoles               = "error adding roles"
	ErrDeletingRoles             = "error deleting roles"
	ErrGettingPermissions        = "error getting permissions"
	ErrAddingPermissions         = "error adding permissions"
	ErrDeletingPermissions       = "error deleting permissions"
	ErrGettingUserPermits        = "error getting user permits"
	ErrUserAlreadyHaveThisPermit = "user already have this permit"
	ErrGettingUser               = "error getting user"
)

func PrintError(msg string, err error) error {
	if err == nil {
		return fmt.Errorf(msg, "")
	}

	return fmt.Errorf(msg+" -> %v", err.Error())
}
