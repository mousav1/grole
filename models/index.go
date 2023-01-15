package models

import "gorm.io/gorm"

type Database struct {
	DB *gorm.DB
}

var connection *Database

func Initializers(db *gorm.DB) *Database {
	connection = &Database{
		DB: db,
	}
	return connection
}
