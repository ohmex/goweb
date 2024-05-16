package migrations

import (
	"goweb/models"

	"gorm.io/gorm"
)

type CreateTables struct{}

func (CreateTables) Id() string {
	return "UserMigration"
}
func (CreateTables) Up(db *gorm.DB) {
	db.Migrator().CreateTable(&models.User{})
	db.Migrator().CreateTable(&models.Post{})
}

func (CreateTables) Down(db *gorm.DB) {
	db.Migrator().DropTable(&models.Post{})
	db.Migrator().DropTable(&models.User{})
}
