package test

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {

	dsn := "host=localhost user=root password=secret dbname=grole port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	os.Exit(m.Run())
}
