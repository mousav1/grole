package models

type Role struct {
	ID          uint `gorm:"primary_key, AUTO_INCREMENT"`
	Name        string
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Role) TableName() string {
	return "roles"
}
