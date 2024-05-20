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
	db.Migrator().AutoMigrate(&models.User{}, &models.Tenant{}, &models.Post{})
}

func (DatabaseTables) Down(db *gorm.DB) {
	db.Migrator().DropTable(&models.Post{})
	db.Migrator().DropTable(&models.Tenant{})
	db.Migrator().DropTable(&models.User{})
}
