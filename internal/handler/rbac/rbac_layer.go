package rbac

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/G0tem/go-service-auth/internal"
	"github.com/G0tem/go-service-auth/internal/model"
	"github.com/G0tem/go-service-auth/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SyncRbacModelAction uint8

const (
	AddMissedOnly SyncRbacModelAction = iota
	DeleteOutdatedRolesAndPermissions
	DeleteLinksBetweenRolesAndPermissions
)

type RBACLayer struct {
	DB  *gorm.DB
	Ctx context.Context
}

// Direct migrations
func (g *RBACLayer) MigrateTables() (err error) {
	if err = g.DB.AutoMigrate(&model.UserRole{}); err != nil {
		return
	}

	if err = g.DB.AutoMigrate(&model.UserPermission{}); err != nil {
		return
	}

	if err = g.DB.AutoMigrate(&model.UserRolePermission{}); err != nil {
		return
	}

	return
}

func (layer *RBACLayer) CheckAccess(rbacList []string) fiber.Handler {
	// Return middleware handler
	return func(fiberCtx *fiber.Ctx) error {
		var (
			userId          uuid.UUID
			userPermits     []string
			permitted       bool
			permits         string
			rolePermissions []string
		)

		err := layer.DB.Table("users").
			Select("role").
			Where("id"+" = ?", userId).
			Find(&permits).Error
		if err != nil {
			log.Error().Err(err)
		}

		err = json.Unmarshal([]byte(permits), &userPermits)
		if err != nil {
			log.Error().Err(err)
		}

		_ = layer.DB.Model(&model.UserRolePermission{}).
			Select("concat(permission.model, ':', permission.action)").
			Joins("FULL JOIN user_roles as role ON role.id = user_role_permissions.role_id").
			Joins("FULL JOIN user_permissions as permission ON permission.id = user_role_permissions.permission_id").
			Where("role.name IN ?", userPermits).
			Find(&rolePermissions)

		userPermits = append(userPermits, rolePermissions...)

		for _, u := range userPermits {
			if internal.StringInSlice(u, rbacList) {
				permitted = true
			}
		}

		if permitted {
			return fiberCtx.Next()
		} else {
			return fiberCtx.SendString(internal.ErrAccessNotPermitted)
		}
	}
}

func (layer *RBACLayer) ValidatePermits(userPermits datatypes.JSON) error {
	roleList := []model.UserRole{}
	layer.DB.Find(&roleList)

	roles := []string{}
	for _, role := range roleList {
		roles = append(roles, role.Name)
	}

	permissionList := []model.UserPermission{}
	layer.DB.Find(&permissionList)

	permissions := []string{}
	for _, permission := range permissionList {
		permissions = append(permissions, permission.Model+":"+permission.Action)
	}

	permits := []string{}
	if len(userPermits) > 0 {
		err := json.Unmarshal(userPermits, &permits)
		if err != nil {
			return internal.PrintError(internal.ErrUnmarshalingPermits, err)
		}
	}

	for _, permit := range permits {
		if !internal.StringInSlice(permit, roles) && !internal.StringInSlice(permit, permissions) {
			return internal.PrintError(internal.ErrPermitNotFound, fmt.Errorf(""))
		}
	}

	return nil
}

// GetPermissions - get stored permissions in the database
func (layer *RBACLayer) GetPermissions() (permissions []model.UserPermission, err error) {
	err = layer.DB.Find(&permissions).Error
	if err != nil {
		return nil, internal.PrintError(internal.ErrGettingPermissions, err)
	}

	return
}

// AddPermission - stores a permission in the database
func (layer *RBACLayer) AddPermission(userPermission *model.UserPermission) (*model.UserPermission, error) {
	err := layer.DB.Where(
		"concat(user_permissions.model,':',user_permissions.action) = ?",
		userPermission.Model+":"+userPermission.Action,
	).FirstOrCreate(userPermission).Error
	if err != nil {
		return &model.UserPermission{}, internal.PrintError(internal.ErrAddingPermissions, err)
	}

	return userPermission, err
}

// DeletePermission - removes a stored permission from database
func (layer *RBACLayer) DeletePermission(permissionModel, permissionAction string) (err error) {
	err = layer.DB.Unscoped().Where("model = ? and action = ?", permissionModel, permissionAction).Delete(&model.UserPermission{}).Error
	if err != nil {
		return internal.PrintError(internal.ErrDeletingPermissions, err)
	}

	return
}

// GetRoles - get stored roles in the database
func (layer *RBACLayer) GetRoles() (roles []model.UserRole, err error) {
	err = layer.DB.Find(&roles).Error
	if err != nil {
		return nil, internal.PrintError(internal.ErrGettingRoles, err)
	}

	return
}

// AddRole - stores a role in the database
func (layer *RBACLayer) AddRole(role model.UserRole) (model.UserRole, error) {
	err := layer.DB.Where("name = ?", role.Name).FirstOrCreate(&role).Error
	if err != nil {
		return model.UserRole{}, internal.PrintError(internal.ErrAddingRoles, err)
	}

	return role, err
}

// DeleteRole  - removes a stored role from database
func (layer *RBACLayer) DeleteRole(roleName string) (err error) {
	err = layer.DB.Unscoped().Where("name = ?", roleName).Delete(&model.UserRole{}).Error
	if err != nil {
		return internal.PrintError(internal.ErrDeletingRoles, err)
	}

	return
}

// GetAllRolePermissions - get stored role permissions in the database
func (layer *RBACLayer) GetAllRolePermissions() (rolePermissions []model.UserRolePermission, err error) {
	err = layer.DB.Preload(clause.Associations).Find(&rolePermissions).Error
	if err != nil {
		return nil, internal.PrintError(internal.ErrGettingRolePermissions, err)
	}

	return
}

// GetRolePermissions - get stored role permissions in the database
func (layer *RBACLayer) GetRolePermissions(roles ...string) (permissions []model.UserPermission, err error) {
	err = layer.DB.Model(&model.UserRolePermission{}).
		Joins("JOIN user_roles ON user_roles.id = user_role_permissions.role_id AND user_roles.name IN ? JOIN user_permissions ON user_permissions.id = user_role_permissions.permission_id", roles).
		Select("user_permissions.*").Find(&permissions).Error
	if err != nil {
		return nil, internal.PrintError(internal.ErrGettingRolePermissions, err)
	}

	return
}

// AddRolePermission - stores a role permission in the database
func (layer *RBACLayer) AddRolePermission(addRolePermission types.AddRolePermission) (
	userRolePermission model.UserRolePermission, err error,
) {
	var (
		role       model.UserRole
		permission model.UserPermission
	)

	existedRoleTx := layer.DB.Where("name = ?", addRolePermission.Role).First(&role)

	existedPermissionTx := layer.DB.Where("model = ? and action = ?",
		addRolePermission.PermissionModel, addRolePermission.PermissionAction,
	).First(&permission)

	if existedRoleTx.RowsAffected == 0 || existedPermissionTx.RowsAffected == 0 {
		return model.UserRolePermission{}, internal.PrintError(internal.ErrRoleOrPermissionNotFound, err)
	}

	userRolePermission.RoleID = role.ID
	userRolePermission.PermissionID = permission.ID

	err = layer.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRolePermission).Error
	if err != nil {
		return model.UserRolePermission{}, internal.PrintError(internal.ErrAddingRolePermission, err)
	}

	return
}

// DeleteRolePermission - removes a stored role permission from database
func (layer *RBACLayer) DeleteRolePermission(roleName, permissionModel, permissionAction string) (err error) {
	var (
		role       model.UserRole
		permission model.UserPermission
	)

	if err = layer.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return
	}

	if err = layer.DB.Where("model = ? and action = ?", permissionModel, permissionAction).First(&permission).Error; err != nil {
		return
	}

	err = layer.DB.Unscoped().Where("role_id = ? and permission_id = ?", role.ID, permission.ID).Delete(&model.UserRolePermission{}).Error
	if err != nil {
		return internal.PrintError(internal.ErrDeletingRolePermission, err)
	}

	return
}

// GetUserPermits - get user roles and permissions in the database
func (layer *RBACLayer) GetUserPermits(userId uuid.UUID) (permits types.Permits, err error) {
	userPermits := &types.UserPermit{}

	err = layer.DB.Table("users").
		Select("id"+" as user_id", "roles"+" as permits").
		Where("id"+" = ?", userId).
		Find(&userPermits).Error
	if err != nil {
		return permits, internal.PrintError(internal.ErrGettingUser, err)
	}

	if userPermits.UserID == "" {
		return permits, internal.PrintError(internal.ErrUserNotFound, err)
	}

	permits.UserID = userPermits.UserID
	permits.Permits = userPermits.Permits

	if len(permits.Permits) > 0 {
		err = json.Unmarshal(permits.Permits, &permits.PermitList)
		if err != nil {
			return permits, internal.PrintError(internal.ErrUnmarshalingPermits, err)
		}
	}

	return
}

// GrantUserRole - grant user role in the database
func (layer *RBACLayer) GrantUserRole(userId uuid.UUID, roleName string) (permits types.Permits, err error) {
	res := layer.DB.Where("name = ?", roleName).First(&model.UserRole{})
	if res.RowsAffected == 0 || res.Error != nil {
		err = res.Error
		return permits, internal.PrintError(internal.ErrRoleNotFound, err)
	}

	return layer.grantUserPermit(userId, roleName)
}

// GrantUserPermission - grant user permission in the database
func (layer *RBACLayer) GrantUserPermission(userId uuid.UUID, permissionModel, permissionAction string) (permits types.Permits, err error) {
	res := layer.DB.Where("model = ? and action = ?", permissionModel, permissionAction).First(&model.UserPermission{})
	if res.RowsAffected == 0 || res.Error != nil {
		err = res.Error
		return permits, internal.PrintError(internal.ErrPermissionNotFound, err)
	}

	return layer.grantUserPermit(userId, permissionModel+":"+permissionAction)
}

// grantUserPermit - grant user role and permission in the database
func (layer *RBACLayer) grantUserPermit(userId uuid.UUID, permit string) (permits types.Permits, err error) {
	permits, err = layer.GetUserPermits(userId)
	if err != nil {
		return permits, internal.PrintError(internal.ErrGettingUserPermits, err)
	}

	if internal.ExistPermit(permits.PermitList, permit) {
		return permits, internal.PrintError(internal.ErrUserAlreadyHaveThisPermit, err)
	}

	permits.PermitList = append(permits.PermitList, permit)

	permits.Permits, err = json.Marshal(permits.PermitList)
	if err != nil {
		return permits, internal.PrintError(internal.ErrMarshallingPermits, err)
	}

	err = layer.DB.Table("users").
		Where("id"+" = ?", userId).
		Updates(map[string]interface{}{"role": permits.Permits}).Error
	if err != nil {
		return permits, internal.PrintError(internal.ErrGrantingPermit, err)
	}

	return
}

// RevokeUserRole - revoke user role in the database
func (layer *RBACLayer) RevokeUserRole(userId uuid.UUID, roleName string) (permits types.Permits, err error) {
	return layer.RevokeUserPermit(userId, roleName)
}

// RevokeUserPermit - revoke user role and permission in the database
func (layer *RBACLayer) RevokeUserPermit(userId uuid.UUID, permit string) (permits types.Permits, err error) {
	var oldPermits []string

	permits, err = layer.GetUserPermits(userId)
	if err != nil {
		return permits, internal.PrintError(internal.ErrGettingUserPermits, err)
	}

	oldPermits = permits.PermitList
	permits.PermitList = []string{}

	for _, p := range oldPermits {
		if p != permit {
			permits.PermitList = append(permits.PermitList, p)
		}
	}

	permits.Permits, err = json.Marshal(permits.PermitList)
	if err != nil {
		return permits, internal.PrintError(internal.ErrMarshallingPermits, err)
	}

	err = layer.DB.Table("users").
		Where("id"+" = ?", userId).
		Updates(map[string]interface{}{"role": permits.Permits}).Error
	if err != nil {
		return permits, internal.PrintError(internal.ErrRevokingPermit, err)
	}

	return
}

func (rbac *RBACLayer) InitSafety(matrix map[string]string) error {
	return rbac.Init(matrix, DeleteLinksBetweenRolesAndPermissions)
}

func (rbac *RBACLayer) Init(matrix map[string]string, action SyncRbacModelAction) error {
	// Get roles, permissions and links between them from storage (database)
	permissions, err := rbac.GetPermissions()
	if err != nil {
		return err
	}
	rolesList, err := rbac.GetRoles()
	if err != nil {
		return err
	}
	allUserRolePermissions, err := rbac.GetAllRolePermissions()
	if err != nil {
		return err
	}
	allUserRolePermissionsIndex := internal.MappingToMultiMap(allUserRolePermissions,
		func(userRolePermission model.UserRolePermission) (string, *model.UserPermission) {
			return userRolePermission.Role.Name, userRolePermission.Permission
		})

	// Collect changes
	var addUserRoleEvents []model.UserRole
	var addPermissionEvents []model.UserPermission
	var addUserRolePermissionEvents []model.UserRolePermission
	var removeUserRoleEvents []uuid.UUID
	var removePermissionEvents []uuid.UUID
	var removeUserRolePermissionEvents [][]uuid.UUID

	rolesBucket := internal.Mapping(rolesList, func(ur model.UserRole) string { return ur.Name })
	permissionsBucket := internal.Mapping(permissions,
		func(userPermission model.UserPermission) string {
			return fmt.Sprintf("%s:%s", userPermission.Model, userPermission.Action)
		})
	for role, rolePermissions := range matrix {
		var userRole model.UserRole
		if !slices.Contains(rolesBucket, role) {
			addUserRoleEvents = append(addUserRoleEvents, model.UserRole{
				Name: role,
			})
			// Hack for prevent double cycling for collect addUserRolePermissionEvents
			userRole, err = rbac.AddRole(model.UserRole{
				Name: role,
			})
			if err != nil {
				return err
			}
			rolesList = append(rolesList, userRole)

		} else {
			userRole = internal.Any(rolesList, func(ur model.UserRole) bool { return ur.Name == role })
			rolesBucket = internal.Remove(rolesBucket, role)
		}

		// Snapshot for current role permissions
		currentRolePermissions, ok := allUserRolePermissionsIndex[role]
		if !ok {
			currentRolePermissions = make([]*model.UserPermission, 0, 10)
		}
		rolePermissionsBucket := internal.Mapping(currentRolePermissions,
			func(x *model.UserPermission) []uuid.UUID {
				return []uuid.UUID{userRole.ID, x.ID}
			})

		for _, permission := range strings.Split(rolePermissions, ",") {
			modelName, action := SplitAndTrim2(permission, ":", " ")
			var userPermission model.UserPermission

			if !slices.Contains(permissionsBucket, permission) {
				userPermission = model.UserPermission{
					Model:  modelName,
					Action: action,
				}
				addPermissionEvents = append(addPermissionEvents, userPermission)
				// Hack for prevent double cycling for collect addUserRolePermissionEvents
				rolePermission, err := rbac.AddPermission(&userPermission)
				if err != nil {
					return err
				}
				userPermission.ID = rolePermission.ID
				permissions = append(permissions, userPermission)

			} else {
				userPermission = internal.Any(permissions, func(x model.UserPermission) bool { return x.Model == modelName && x.Action == action })
				permissionsBucket = internal.Remove(permissionsBucket, permission)
			}

			// Collect links between roles and permissions
			if slices.IndexFunc(currentRolePermissions, func(x *model.UserPermission) bool { return x.ID == userPermission.ID }) < 0 {
				addUserRolePermissionEvents = append(addUserRolePermissionEvents, model.UserRolePermission{
					RoleID:       userRole.ID,
					PermissionID: userPermission.ID,
				})
			} else {
				rolePermissionsBucket = internal.RemoveFunc(rolePermissionsBucket, func(x []uuid.UUID) bool { return x[0] == userRole.ID && x[1] == userPermission.ID })
			}
		}

		// Hack for prevent double cycling for collect addUserRolePermissionEvents
		removeUserRolePermissionEvents = append(removeUserRolePermissionEvents, rolePermissionsBucket...)
	}
	for _, roleName := range rolesBucket {
		userRole := internal.Any(rolesList, func(x model.UserRole) bool { return x.Name == roleName })
		removeUserRoleEvents = append(removeUserRoleEvents, userRole.ID)
	}
	for _, permission := range permissionsBucket {
		modelName, action := SplitAndTrim2(permission, ":", " ")
		userPermission := internal.Any(permissions, func(x model.UserPermission) bool {
			return x.Model == modelName && x.Action == action
		})
		removePermissionEvents = append(removePermissionEvents, userPermission.ID)
	}

	// Process events
	// For add events only links between roles and permissions process cause hacking above
	for _, userRolePermission := range addUserRolePermissionEvents {
		err = rbac.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRolePermission).Error
		if err != nil {
			return internal.PrintError(internal.ErrAddingRolePermission, err)
		}
	}

	// Optionally process remove events
	switch action {
	case DeleteLinksBetweenRolesAndPermissions:

		for _, removeEvent := range removeUserRolePermissionEvents {
			var userRolePermission model.UserRolePermission
			err = rbac.DB.Where(&model.UserRolePermission{
				RoleID:       removeEvent[0],
				PermissionID: removeEvent[1],
			}).Delete(&userRolePermission).Error
			if err != nil {
				return internal.PrintError(internal.ErrDeletingRolePermission, err)
			}

		}

	case DeleteOutdatedRolesAndPermissions:

		for _, removeEvent := range removeUserRolePermissionEvents {
			var userRolePermission model.UserRolePermission
			err = rbac.DB.Where(&model.UserRolePermission{
				RoleID:       removeEvent[0],
				PermissionID: removeEvent[1],
			}).Delete(&userRolePermission).Error
			if err != nil {
				return internal.PrintError(internal.ErrDeletingRolePermission, err)
			}

		}

		for _, removeEvent := range removePermissionEvents {
			var userPermission model.UserPermission
			err = rbac.DB.Where(&model.UserPermission{
				ID: removeEvent,
			}).Delete(&userPermission).Error
			if err != nil {
				return internal.PrintError(internal.ErrDeletingPermissions, err)
			}

		}

		for _, removeEvent := range removeUserRoleEvents {
			var userRole model.UserRole
			err = rbac.DB.Where(&model.UserRole{
				ID: removeEvent,
			}).Delete(&userRole).Error
			if err != nil {
				return internal.PrintError(internal.ErrDeletingRoles, err)
			}

		}
	}

	log.Debug().Msgf("Init rbac layer stat: added roles: %v, added permissions %v", len(addUserRoleEvents), len(addPermissionEvents))

	return nil
}

func Split2(text string, separator string) (string, string) {
	result := strings.SplitN(text, separator, 2)
	return result[0], result[1]
}

func SplitAndTrim2(text string, separator string, cutset string) (string, string) {
	result := strings.SplitN(text, separator, 2)
	return strings.Trim(result[0], cutset), strings.Trim(result[1], cutset)
}

func SplitAndTrim(text string, separator string, cutset string) []string {
	result := strings.Split(text, separator)
	return internal.Mapping(result, func(x string) string { return strings.Trim(x, cutset) })
}
