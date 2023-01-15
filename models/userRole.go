package models

type UserRoles struct {
	UserID uint `gorm:"primaryKey" column:"user_id"`
	RoleID uint `gorm:"primaryKey" column:"role_id"`
}

func (UserRoles) TableName() string {
	return "user_roles"
}
