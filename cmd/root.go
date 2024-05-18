package cmd

import (
	"fmt"
	"goweb/config"
	"goweb/docs"
	"goweb/routes"
	"goweb/server"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goweb",
	Short: "Go web application",
	Long: `An application showcasing the following features:
1. REST Api implementation
2. Cobra commands
3. GORM implementation with migrations`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Interface("Port", cmd.Flag("port").Value).Send()
		StartServer()
	},
}

func init() {
	rootCmd.Flags().Int16P("port", "p", 8080, "Port to run server on.")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func StartServer() {
	cfg := config.NewConfig()
	// TODO: overide from cmdline options
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	app := server.NewServer(cfg)
	routes.ConfigureRoutes(app)
	err := app.Start(cfg.HTTP.Port)
	if err != nil {
		log.Fatal().Msg("Port already used")
	}
}
