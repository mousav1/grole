package migrate

import (
	"github.com/mousav1/grole/models"
	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB) {
	db.AutoMigrate(&models.Permission{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.UserRoles{})
}
