package migrations

import (
	"goweb/models"
	"os"

	"gorm.io/gorm"
)

type DatabaseTables struct{}

func (DatabaseTables) Id() string {
	return "UserMigration"
}
func (DatabaseTables) Up(db *gorm.DB) {
	db.SetupJoinTable(&models.Domain{}, "Users", &models.DomainUser{})
	db.Migrator().AutoMigrate(&models.Role{}, &models.User{}, &models.Domain{})

	dbType := db.Dialector.Name()
	partitioningEnabled := false
	if os.Getenv("DB_PARTITIONING_ENABLED") == "true" {
		partitioningEnabled = true
	}

	if dbType == "postgres" && partitioningEnabled {
		db.Exec(`
			CREATE TABLE posts (
				id BIGSERIAL,
				uuid CHAR(36),
				created_at TIMESTAMP,
				updated_at TIMESTAMP,
				deleted_at TIMESTAMP,
				domain CHAR(36) NOT NULL,
				title TEXT,
				content TEXT,
				user_id INTEGER,
				PRIMARY KEY (id, domain),
				UNIQUE (uuid, domain)
			) PARTITION BY LIST (domain);
		`)
		// Optionally, create a default partition or per-domain partitions
		db.Exec(`CREATE TABLE posts_default PARTITION OF posts DEFAULT;`)
	}

	// Let GORM add any missing columns/indexes
	db.Migrator().AutoMigrate(&models.Post{})

}

func (DatabaseTables) Down(db *gorm.DB) {
	db.Migrator().DropTable(&models.Post{})
	db.Migrator().DropTable(&models.Domain{})
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Role{})
}
