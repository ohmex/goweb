package main

import (
	"gowebmvc/server"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	key := os.Getenv("KEY")
	log.Info().Str("KEY", key).Send()

	if err := godotenv.Load(); err != nil {
		log.Error().Str("ERROR", "error loading .env file")
	}
}

func main() {
	app := server.New()

	//data, _ := json.MarshalIndent(e.Routes(), "", "  ")
	//log.Printf("\n%s", data)

	app.Logger.Fatal(app.Start(":8080"))
}
