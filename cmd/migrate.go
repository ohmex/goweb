package cmd

import (
	"goweb/config"
	"goweb/db"
	"goweb/db/migrations"
	"goweb/services"
	"goweb/util"

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

	// Create partitions for existing domains after migration
	CreatePartitionsForExistingDomains()
}

func MigrateDown() {
	db.AddMigrators(migrations.DatabaseTables{}, migrations.TableData{})

	if err := db.MigrateDown(GetDB()); err != nil {
		log.Fatal().Msg("Migrate DOWN failed")
	}
}

func CreatePartitionsForExistingDomains() {
	database := GetDB()
	domainService := services.NewDomainService(database)

	// Check if partitioning is enabled
	if !util.IsPartitioningEnabled() {
		log.Info().Msg("Partitioning disabled - skipping partition creation")
		return
	}

	// Check if database supports partitioning
	if !util.IsDatabasePartitioningSupported(database) {
		log.Info().Msg("Database does not support partitioning - skipping partition creation")
		return
	}

	if err := domainService.CreatePartitionsForExistingDomains(); err != nil {
		log.Fatal().Err(err).Msg("Failed to create partitions for existing domains")
	}

	log.Info().Msg("Successfully created partitions for all existing domains")
}
