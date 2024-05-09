package main

import (
	"gowebmvc/framework"
	"gowebmvc/model"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("KEY")
	log.Printf("\n>>>> KEY: %s", key)

	e := framework.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Resource("/user", model.User{})

	//data, _ := json.MarshalIndent(e.Routes(), "", "  ")
	//log.Printf("\n%s", data)

	e.Logger.Fatal(e.Start(":8080"))
}
