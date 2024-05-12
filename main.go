package main

import (
	"gowebmvc/server"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	if err := godotenv.Load(); err != nil {
		log.Error().Str("ERROR", "error loading .env file")
	}
}

func main() {
	app := server.New()

	log.Info().Interface("Routes", app.Routes()).Send()

	app.Logger.Fatal(app.Start(":8080"))
}
