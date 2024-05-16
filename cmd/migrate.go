package cmd

import (
	"goweb/config"
	"goweb/db"
	"goweb/db/migrations"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
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
	cfg := config.NewConfig()
	return db.InitDB(cfg)
}

func MigrateUp() {
	db.AddMigrators(migrations.DatabaseTables{}, migrations.TableData{})

	if err := db.Migrate(GetDB()); err != nil {
		log.Fatal().Msg("Migrate UP failed")
	}
}

func MigrateDown() {
	db.AddMigrators(migrations.DatabaseTables{}, migrations.TableData{})

	if err := db.MigrateDown(GetDB()); err != nil {
		log.Fatal().Msg("Migrate DOWN failed")
	}
}
