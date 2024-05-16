package cmd

import (
	"goweb/db"
	"goweb/db/migrations"
	"os"

	"github.com/rs/zerolog/log"

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

func GetDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"
	gorm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return gorm
}

func MigrateUp() {
	db.AddMigrators(migrations.CreateTables{}, migrations.InsertData{})

	if err := db.Migrate(GetDB()); err != nil {
		log.Fatal().Msg("Migrate UP failed")
	}
}

func MigrateDown() {
	db.AddMigrators(migrations.CreateTables{}, migrations.InsertData{})

	if err := db.MigrateDown(GetDB()); err != nil {
		log.Fatal().Msg("Migrate DOWN failed")
	}
}
