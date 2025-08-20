package model

import "github.com/google/uuid"

// RolePermission represents the database that stores the relationship between roles and permissions
type UserRolePermission struct {
	RoleID       uuid.UUID       `json:"-" gorm:"primaryKey;type:uuid"`
	Role         *UserRole       `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PermissionID uuid.UUID       `json:"-" gorm:"primaryKey;type:uuid"`
	Permission   *UserPermission `gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
