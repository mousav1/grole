package models

type Permission struct {
	ID          uint `gorm:"primary_key, AUTO_INCREMENT"`
	Name        string
	Description string
	Roles       []Role `gorm:"many2many:permission_role;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Permission) TableName() string {
	return "permissions"
}
