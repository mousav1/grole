package grole

import (
	"errors"

	"github.com/mousav1/grole/migrate"
	"github.com/mousav1/grole/models"
	"gorm.io/gorm"
)

type Options struct {
	DB *gorm.DB
}

var conn *models.Database

// set database connection
func New(opt Options) *models.Database {
	conn = models.Initializers(opt.DB)
	migrate.MigrateTables(opt.DB)
	return conn
}

// delete the given role
// @param uint
// @return bool, error
func DeleteRole(roleId uint) (bool, error) {

	var userRole models.UserRoles
	res := conn.DB.Where("role_id = ?", roleId).First(&userRole)
	if res.Error == nil {
		return false, errors.New("ROLE IS ASSIGNED")
	}

	res = conn.DB.Where("id = ?", roleId).Delete(&models.Role{})
	if res.Error != nil {
		return false, res.Error
	} else if res.RowsAffected < 1 {
		return false, errors.New("CANNOT BE DELETED BECAUSE IT DOESN'T EXIST")
	}
	return true, nil
}

// update the given role
// @param uint
// @return bool, error
func UpdateRole(roleId uint, newRole models.Role) (bool, error) {
	res := conn.DB.Where("id = ?", roleId).Updates(models.Role{Name: newRole.Name, Description: newRole.Description})
	if res.Error != nil {
		return false, res.Error
	} else if res.RowsAffected < 1 {
		return false, errors.New("CANNOT BE UPDATE BECAUSE IT DOESN'T EXIST")
	}
	return true, nil
}

// delete the given permission
// @param uint
// @return bool, error
func DeletePermission(permissionId uint) (bool, error) {
	var permission models.Permission

	res := conn.DB.Where("id = ?", permissionId).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
	}

	roleCount := conn.DB.Model(&permission).Association("Roles").Count()
	if roleCount > 0 {
		return false, errors.New("PERMISSION IS ASSIGNED")
	}

	res = conn.DB.Where("id = ?", permissionId).Delete(&models.Permission{})
	if res.Error != nil {
		return false, res.Error
	} else if res.RowsAffected < 1 {
		return false, errors.New("CANNOT BE DELETED BECAUSE IT DOESN'T EXIST")
	}

	return true, nil
}

// update the given permission
// @param uint
// @return bool, error
func UpdatePermission(permissionId uint, newRole models.Permission) (bool, error) {
	res := conn.DB.Where("id = ?", permissionId).Updates(models.Permission{Name: newRole.Name, Description: newRole.Description})
	if res.Error != nil {
		return false, res.Error
	} else if res.RowsAffected < 1 {
		return false, errors.New("CANNOT BE UPDATE BECAUSE IT DOESN'T EXIST")
	}
	return true, nil
}

// Show All Permission With Role
// @return []models.Permission, error
func FindAllPermission() ([]models.Permission, error) {
	var permissions []models.Permission
	res := conn.DB.Preload("Roles").Find(&permissions)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("RECORD NOT FOUND")
		} else {
			return nil, res.Error
		}
	}
	return permissions, nil
}

// Find roles of each permission
// @param string
// @return []models.Role
func Roles(permissions ...string) ([]models.Role, error) {
	var roles []models.Role
	var allRole []models.Role
	for _, permission := range permissions {
		per, err := FindPermissionByName(permission)
		if err != nil {
			return nil, errors.New("PERMISSION DOESN'T EXIST")
		}
		conn.DB.Model(&per).Association("Roles").Find(&roles)
		allRole = append(allRole, roles...)
	}
	return allRole, nil
}

// find Permission By Name and Show each with Role
// @param string
// @return models.Permission, error
func FindPermissionByName(name string) (models.Permission, error) {
	var permission models.Permission
	res := conn.DB.Where("name = ?", name).Preload("Roles").First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return permission, errors.New("PERMISSION NOT FOUND")
		}
	}
	return permission, nil
}

// find Permission By Id and Show each with Role
// @param uint
// @return models.Permission, error
func FindPermissionById(id uint) (models.Permission, error) {
	var permission models.Permission
	res := conn.DB.Where("id = ?", id).Preload("Roles").First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return permission, errors.New("RECORD NOT FOUND")
		}
		return permission, res.Error
	}
	return permission, nil
}

// find Permission or Create Permission If not found
// @param models.Permission
// @return models.Permission, error
func FindOrCreatePermission(permission models.Permission) (models.Permission, error) {
	var newPermission models.Permission
	res := conn.DB.FirstOrCreate(&newPermission, permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return newPermission, errors.New("RECORD NOT FOUND")
		}
		if res.RowsAffected < 1 {
			return newPermission, errors.New("CANNOT BE INSERT OR CREATE PERMISSION ")
		}
		return newPermission, res.Error
	}
	return newPermission, nil
}

// Revoke the given role by id for permission
// @param uint, uint
// @return bool, error
func RemoveRoleByIdFromPermission(permissionId uint, roleId uint) (bool, error) {
	var role models.Role
	var permission models.Permission

	res := conn.DB.Where("id = ?", roleId).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("ROLE NOT FOUND")
		}
	}
	res = conn.DB.Where("id = ?", permissionId).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("PERMISSION NOT FOUND")
		}
	}

	error := conn.DB.Model(&permission).Association("Roles").Delete(&role)
	if error != nil {
		return false, error
	}
	return true, nil
}

// Revoke the given role by name for permission
// @param uint, string
// @return bool, error
func RemoveRoleByNameFromPermission(permissionId uint, roleName string) (bool, error) {
	var role models.Role
	var permission models.Permission

	res := conn.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("ROLE NOT FOUND")
		}
	}
	res = conn.DB.Where("id = ?", permissionId).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("PERMISSSION NOT FOUND")
		}
	}

	error := conn.DB.Model(&permission).Association("Roles").Delete(&role)
	if error != nil {
		return false, error
	}
	return true, nil
}

// Revoke the given role for permission
// @param uint, string
// @return bool, error
func RemoveRoleFromPermission(permissionId uint, roleName string) (bool, error) {
	var permission models.Permission

	res := conn.DB.Where("id = ?", permissionId).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
	}

	findRole, err := FindRoleByName(roleName)
	if err != nil {
		return false, errors.New("ROLE DOESN'T EXIST")
	}

	error := conn.DB.Model(&permission).Association("Roles").Delete(&findRole)
	if error != nil {
		return false, error
	}
	return true, nil
}

// Return the number of permissions Role.
// @param uint
// @return int64, error
func CountRoleFromPermission(permissionId uint) (int64, error) {
	var permission models.Permission

	res := conn.DB.Where("id = ?", permissionId).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("RECORD NOT FOUND")
		}
	}

	return conn.DB.Model(&permission).Association("Roles").Count(), nil
}

// Remove all current Role for Permission.
// @param uint
// @return bool, error
func RemoveAllRoleFromPermission(permissionId uint) (bool, error) {
	var permission models.Permission

	res := conn.DB.Where("id = ?", permissionId).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
	}

	error := conn.DB.Model(&permission).Association("Roles").Clear()
	if error != nil {
		return false, error
	}

	return true, nil
}

// Remove all current Permission role and set the given ones.
// @param uint, string
// @return []models.Role, error
func SyncRolesFromPermission(permissionId uint, roles ...string) ([]models.Role, error) {
	var permission models.Permission
	rolesModel := []models.Role{}

	res := conn.DB.Where("id = ?", permissionId).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("RECORD NOT FOUND")
		}
	}

	for _, roleName := range roles {
		role, err := FindRoleByName(roleName)
		if err != nil {
			return rolesModel, errors.New("ROLE DOESN'T EXIST")
		} else {
			rolesModel = append(rolesModel, role)
		}
	}

	error := conn.DB.Model(&permission).Association("Roles").Replace(rolesModel)
	if error != nil {
		return nil, error
	}
	return rolesModel, nil
}

// Find All Role
// @return []models.Role, error
func FindAllRole() ([]models.Role, error) {
	var roles []models.Role
	res := conn.DB.Preload("Permissions").Find(&roles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("RECORD NOT FOUND")
		}
		return nil, res.Error
	}
	return roles, nil
}

// Return all Permissions the Role.
// @param string
// @return []models.Permission
func Permissions(roles ...string) ([]models.Permission, error) {
	var permissions []models.Permission
	var allPermission []models.Permission
	for _, role := range roles {
		_, err := FindRoleByName(role)
		if err != nil {
			return nil, errors.New("ROLE DOESN'T EXIST")
		}
		conn.DB.Model(&role).Association("Permissions").Find(&permissions)
		allPermission = append(allPermission, permissions...)
	}
	return allPermission, nil
}

// Find Role By Name
// @param string
// @return models.Role, error
func FindRoleByName(name string) (models.Role, error) {
	var role models.Role
	res := conn.DB.Where("name = ?", name).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return role, errors.New("RECORD NOT FOUND")
		}
		return role, res.Error
	}
	return role, nil
}

// Find Role By Id
// @param uint
// @return models.Role, error
func FindRoleById(roleId uint) (models.Role, error) {
	var role models.Role
	res := conn.DB.Where("id = ?", roleId).Preload("Permissions").First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return role, errors.New("RECORD NOT FOUND")
		}
		return role, res.Error
	}
	return role, nil
}

// Find Or Create Role
// @param models.Role
// @return models.Role, error
func FindOrCreateRole(role models.Role) (models.Role, error) {
	var newRole models.Role
	res := conn.DB.FirstOrCreate(&newRole, role)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return newRole, errors.New("RECORD NOT FOUND")
		}
		if res.RowsAffected < 1 {
			return newRole, errors.New("CANNOT BE INSERT OR CREATE DATA ")
		}
		return newRole, res.Error
	}
	return newRole, nil
}

// Get Name Roles
// @param []models.Role
// @return []string
func GetNameRoles(roles []models.Role) []string {
	var rolesName []string
	for _, value := range roles {
		rolesName = append(rolesName, value.Name)
	}
	return rolesName
}

// Revoke the given Permission by id for Role
// @param uint, uint
// @return bool, error
func RemovePermissionByIdFromRole(roleId uint, permissionId uint) (bool, error) {
	var role models.Role
	var permission models.Permission

	res := conn.DB.Where("id = ?", roleId).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
		return false, res.Error
	}
	res = conn.DB.Where("id = ?", permissionId).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
		return false, res.Error
	}

	error := conn.DB.Model(&role).Association("Permissions").Delete(&permission)
	if error != nil {
		return false, error
	}
	return true, nil
}

// Revoke the given Permission by name for Role
// @param string, string
// @return bool, error
func RemovePermissionByNameFromRole(roleName string, permissionName string) (bool, error) {
	var role models.Role
	var permission models.Permission

	res := conn.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
		return false, res.Error
	}
	res = conn.DB.Where("name = ?", permissionName).First(&permission)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
		return false, res.Error
	}

	error := conn.DB.Model(&role).Association("Permissions").Delete(&permission)
	if error != nil {
		return false, error
	}
	return true, nil
}

// Return the number of Role permissions.
// @param uint
// @return int64, error
func CountPermissionFromRole(roleId uint) (int64, error) {
	var role models.Role

	res := conn.DB.Where("id = ?", roleId).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("RECORD NOT FOUND")
		}
		return 0, res.Error
	}
	return conn.DB.Model(&role).Association("Permissions").Count(), nil
}

// Remove all current Permission for Role.
// @param uint
// @return bool, error
func RemoveAllPermissionFromRole(roleId uint) (bool, error) {
	var role models.Role

	res := conn.DB.Where("id = ?", roleId).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
		return false, res.Error
	}

	error := conn.DB.Model(&role).Association("Permissions").Clear()
	if error != nil {
		return false, error
	}

	return true, nil
}

// Remove all current role Permission and set the given ones.
// @param uint, string
// @return []models.Permission, error
func SyncPermissionsFromRole(roleId uint, permissions ...string) ([]models.Permission, error) {
	var role models.Role
	permissionModels := []models.Permission{}

	res := conn.DB.Where("id = ?", roleId).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("RECORD NOT FOUND")
		}
		return nil, res.Error
	}

	for _, permissionName := range permissions {
		permission, err := FindPermissionByName(permissionName)
		if err != nil {
			return permissionModels, errors.New("PERMISSION DOESN'T EXIST")
		} else {
			permissionModels = append(permissionModels, permission)
		}
	}

	error := conn.DB.Model(&role).Association("Permissions").Replace(&permissionModels)
	if error != nil {
		return nil, error
	}
	return permissionModels, nil
}

// Assign the given Permissions to the Role.
// @param uint, string
// @return []models.Permission, error
func AssignPermissionsFromRole(roleId uint, permissions ...string) ([]models.Permission, error) {
	var role models.Role
	permissionModels := []models.Permission{}

	res := conn.DB.Where("id = ?", roleId).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("RECORD NOT FOUND")
		}
		return nil, res.Error
	}

	for _, permissionName := range permissions {
		permission, err := FindPermissionByName(permissionName)
		if err != nil {
			return permissionModels, errors.New("PERMISSION DOESN'T EXIST")
		} else {
			permissionModels = append(permissionModels, permission)
		}
	}

	error := conn.DB.Model(&role).Association("Permissions").Append(&permissionModels)
	if error != nil {
		return nil, error
	}
	return permissionModels, nil
}

// Determine if the Role may perform the given permission.
// @param uint, string
// @return models.Permission, error
func HasPermissionTo(roleId uint, permissionName string) (models.Permission, error) {
	var role models.Role

	res := conn.DB.Where("id = ?", roleId).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return models.Permission{}, errors.New("RECORD NOT FOUND")
		}
		return models.Permission{}, res.Error
	}
	var permission models.Permission

	permissionId, error := FindPermissionByName(permissionName)
	if error != nil {
		return permission, error
	}

	error = conn.DB.Model(&role).Where("permission_id = ?", permissionId.ID).Association("Permissions").Find(&permission)
	if error != nil {
		return permission, error
	}
	return permission, nil
}

// Return all the Roles the user.
// @param uint
// @return []models.Role, error
func GetRole(userID uint) ([]models.Role, error) {
	var userRoles []models.UserRoles
	res := conn.DB.Where("user_id = ?", userID).Find(&userRoles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []models.Role{}, errors.New("RECORD NOT FOUND")
		}
		return []models.Role{}, res.Error
	}

	var roles []models.Role
	for _, r := range userRoles {
		var role models.Role
		res := conn.DB.Where("id = ?", r.RoleID).Find(&role)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return []models.Role{}, errors.New("ROLE ID NOT FOUND")
			}
			return []models.Role{}, res.Error
		}
		if res.Error == nil {
			roles = append(roles, role)
		}
	}
	return roles, nil
}

// Return all the Roles Name the user.
// @param uint
// @return []string, error
func GetRoleNames(userID uint) ([]string, error) {
	var userRoles []models.UserRoles
	res := conn.DB.Where("user_id = ?", userID).Find(&userRoles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("RECORD NOT FOUND")
		}
		return nil, res.Error
	}

	var roles []models.Role
	for _, r := range userRoles {
		var role models.Role
		res := conn.DB.Where("id = ?", r.RoleID).Find(&role)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return nil, errors.New("ROLE ID NOT FOUND")
			}
			return nil, res.Error
		}
		if res.Error == nil {
			roles = append(roles, role)
		}
	}
	return GetNameRoles(roles), nil
}

// Return all the permissions the user.
// @param uint
// @return []models.Permission, error
func GetAllPermissions(userID uint) ([]models.Permission, error) {
	var userRoles []models.UserRoles
	res := conn.DB.Where("user_id = ?", userID).Find(&userRoles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("RECORD NOT FOUND")
		}
		return nil, res.Error
	}

	var roles []string
	for _, r := range userRoles {
		var role models.Role
		res := conn.DB.Where("id = ?", r.RoleID).Find(&role)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return nil, errors.New("ROLE ID NOT FOUND")
			}
			return nil, res.Error
		}
		if res.Error == nil {
			roles = append(roles, role.Name)
		}
	}
	permissions, error := Permissions(roles...)
	if error != nil {
		return nil, error
	}
	return permissions, nil
}

// Assign the given roles to the User.
// @param uint, string
// @return bool, error
func AssignRoles(userID uint, Roles ...string) (bool, error) {
	var userRole models.UserRoles

	for _, roleName := range Roles {
		role, err := FindRoleByName(roleName)
		if err != nil {
			return false, errors.New("ROLE DOESN'T EXIST")
		}
		conn.DB.FirstOrCreate(&userRole, models.UserRoles{
			UserID: userID,
			RoleID: role.ID,
		})
	}

	return true, nil
}

// Revoke the given role by id for user
// @param uint, uint
// @return bool, error
func RemoveRoleByIdFromUser(userID uint, roleId uint) (bool, error) {
	res := conn.DB.Where("user_id = ?", userID).Where("role_id = ?", roleId).Delete(&models.UserRoles{})
	if res.Error != nil {
		return false, res.Error
	} else if res.RowsAffected < 1 {
		return false, errors.New("CANNOT BE DELETED BECAUSE IT DOESN'T EXIST")
	}
	return true, nil
}

// Revoke the given role by name for user
// @param uint, string
// @return bool, error
func RemoveRoleByNameFromUser(userID uint, roleName string) (bool, error) {
	role, error := FindRoleByName(roleName)
	if error != nil {
		return false, error
	}
	res := conn.DB.Where("user_id = ?", userID).Where("role_id = ?", role.ID).Delete(&models.UserRoles{})
	if res.Error != nil {
		return false, res.Error
	} else if res.RowsAffected < 1 {
		return false, errors.New("CANNOT BE DELETED BECAUSE IT DOESN'T EXIST")
	}
	return true, nil
}

// Remove all current roles for user.
// @param uint
// @return bool, error
func RemoveAllRoleFromUser(userID uint) (bool, error) {
	res := conn.DB.Where("user_id = ?", userID).Delete(&models.UserRoles{})
	if res.Error != nil {
		return false, res.Error
	} else if res.RowsAffected < 1 {
		return false, errors.New("CANNOT BE DELETED BECAUSE IT DOESN'T EXIST")
	}
	return true, nil
}

// Remove all current user roles and set the given ones.
// @param uint, string
// @return bool, error
func SyncRolesFromUser(userID uint, Roles ...string) (bool, error) {
	res := conn.DB.Where("user_id = ?", userID).Delete(&models.UserRoles{})
	if res.Error != nil {
		return false, res.Error
	} else if res.RowsAffected < 1 {
		return false, errors.New("CANNOT BE DELETED BECAUSE IT DOESN'T EXIST")
	}

	for _, roleName := range Roles {
		role, err := FindRoleByName(roleName)
		if err != nil {
			return false, errors.New("ROLE DOESN'T EXIST")
		}
		conn.DB.Create(&models.UserRoles{
			UserID: userID,
			RoleID: role.ID,
		})
	}
	return true, nil
}

// Determine if the user has  of the given role id.
// @param uint, uint
// @return bool, error
func HasRole(userID uint, roleId uint) (bool, error) {
	_, error := FindRoleById(roleId)
	if error != nil {
		return false, error
	}

	var userRole models.UserRoles
	res := conn.DB.Where("user_id = ?", userID).Where("role_id = ?", roleId).First(&userRole)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("RECORD NOT FOUND")
		}
		return false, res.Error
	}
	return true, nil
}

// Determine if the user has of the given roles name.
// @param uint, uint
// @return bool, error
func HasAnyRole(userID uint, rolesName ...string) (bool, error) {
	roles, error := GetRole(userID)
	if error != nil {
		return false, error
	}
	for _, role := range roles {
		for _, name := range rolesName {
			if role.Name == name {
				return true, nil
			}
		}
	}
	return false, nil
}

// Determine if the user has all of the given roles name.
// @param uint, uint
// @return bool, error
func HasAllRole(userID uint, rolesName ...string) (bool, error) {
	roles, error := GetRole(userID)
	if error != nil {
		return false, error
	}

	for index, role := range roles {
		if role.Name != rolesName[index] {
			return false, nil
		}
	}

	return true, nil
}

// Determine if the user has all of the given permissions name.
// @param uint, uint
// @return bool, error
func HasAllPermission(userID uint, permissionsName ...string) (bool, error) {
	permissions, error := GetAllPermissions(userID)
	if error != nil {
		return false, error
	}

	for index, permission := range permissions {
		if permission.Name != permissionsName[index] {
			return false, nil
		}
	}

	return true, nil
}

// Determine if the User has of the given permissions name.
// @param uint, string
// @return bool, error
func HasAnyPermissions(userID uint, permissionsName ...string) (bool, error) {
	permissions, error := GetAllPermissions(userID)
	if error != nil {
		return false, error
	}
	for _, permission := range permissions {
		for _, name := range permissionsName {
			if permission.Name == name {
				return true, nil
			}
		}
	}
	return false, nil
}
