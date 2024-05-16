package migrations

import (
	"goweb/models"

	"gorm.io/gorm"
)

type DatabaseTables struct{}

func (DatabaseTables) Id() string {
	return "UserMigration"
}
func (DatabaseTables) Up(db *gorm.DB) {
	db.Migrator().CreateTable(&models.User{})
	db.Migrator().CreateTable(&models.Post{})
}

func (DatabaseTables) Down(db *gorm.DB) {
	db.Migrator().DropTable(&models.Post{})
	db.Migrator().DropTable(&models.User{})
}
