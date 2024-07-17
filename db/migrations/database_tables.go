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
	db.SetupJoinTable(&models.Domain{}, "Users", &models.DomainUser{})
	db.Migrator().AutoMigrate(&models.Role{}, &models.User{}, &models.Domain{}, &models.Post{})
}

func (DatabaseTables) Down(db *gorm.DB) {
	db.Migrator().DropTable(&models.Post{})
	db.Migrator().DropTable(&models.Domain{})
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Role{})
}
