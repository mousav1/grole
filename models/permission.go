package models

type Permission struct {
	ID          uint `gorm:"primary_key, AUTO_INCREMENT"`
	Name        string
	Description string
	Roles       []Role `gorm:"many2many:permission_role"`
}

func (Permission) TableName() string {
	return "permissions"
}
