package migrations

import (
	"goweb/models"
	"goweb/util"

	"gorm.io/gorm"
)

type DatabaseTables struct{}

func (DatabaseTables) Id() string {
	return "UserMigration"
}
func (DatabaseTables) Up(db *gorm.DB) {
	db.SetupJoinTable(&models.Domain{}, "Users", &models.DomainUser{})
	db.Migrator().AutoMigrate(&models.Role{}, &models.User{}, &models.Domain{})

	partitioningEnabled := util.IsPartitioningEnabled()

	if util.IsDatabasePartitioningSupported(db) && partitioningEnabled {
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

	// Add performance indexes for frequently queried columns
	addPerformanceIndexes(db)
}

func (DatabaseTables) Down(db *gorm.DB) {
	db.Migrator().DropTable(&models.Post{})
	db.Migrator().DropTable(&models.Domain{})
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Role{})
}

// Add performance indexes for frequently queried columns
func addPerformanceIndexes(db *gorm.DB) {
	// User table indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_uuid ON users(uuid)")

	// Domain table indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_domains_uuid ON domains(uuid)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_domains_name ON domains(name)")

	// DomainUser table indexes for join queries
	db.Exec("CREATE INDEX IF NOT EXISTS idx_domain_users_user_id ON domain_users(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_domain_users_domain_id ON domain_users(domain_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_domain_users_active ON domain_users(active)")

	// Post table indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_uuid ON posts(uuid)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_domain_id ON posts(domain_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at)")

	// Role table indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_roles_uuid ON roles(uuid)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_roles_domain_id ON roles(domain_id)")
}
