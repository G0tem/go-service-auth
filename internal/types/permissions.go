package types

import "gorm.io/datatypes"

type UserPermit struct {
	UserID  string         `json:"user_id" gorm:"column:user_id"`
	Permits datatypes.JSON `json:"permits" gorm:"column:permits"`
}

type Permits struct {
	UserID     string         `json:"user_id" gorm:"column:user_id"`
	Name       string         `json:"name" gorm:"column:name"`
	Model      string         `json:"model" gorm:"column:model"`
	Action     string         `json:"action" gorm:"column:action"`
	Permits    datatypes.JSON `json:"permits" gorm:"column:permits"`
	PermitList []string       `json:"permit_list" gorm:"column:permit_list"`
}

// AddRolePermission represents a model used to store RolePermission
type AddRolePermission struct {
	Role             string
	PermissionModel  string
	PermissionAction string
}

type AcceptRoleResponse struct {
	Role string
}
