package cmd

import (
	"goweb/db"
	"goweb/db/migrations"
	"log"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var down bool

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate GORM based database migrations",
	Long:  `Apply & Remove SQL migrations from the applications`,
	Run: func(cmd *cobra.Command, args []string) {
		if down {
			MigrateDown()
		} else {
			MigrateUp()
		}
	},
}

func init() {
	migrateCmd.Flags().BoolVarP(&down, "down", "d", false, "Remove/Destroy migrations")
	rootCmd.AddCommand(migrateCmd)
}

func MigrateUp() {
	db.AddMigrators(migrations.CreateTables{}, migrations.InsertData{})

	dsn := "root:sachin@tcp(127.0.0.1:3306)/gowebmvc?charset=utf8mb4&parseTime=True&loc=Local"
	gorm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.Migrate(gorm); err != nil {
		log.Fatal("Migrate UP failed")
	}
}

func MigrateDown() {
	db.AddMigrators(migrations.CreateTables{}, migrations.InsertData{})

	dsn := "root:sachin@tcp(127.0.0.1:3306)/gowebmvc?charset=utf8mb4&parseTime=True&loc=Local"
	gorm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.Migrate(gorm); err != nil {
		log.Fatal("Migrate DOWN failed")
	}
}
