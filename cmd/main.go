package main

import (
	"fmt"
	application "goweb"
	"goweb/config"
	"goweb/docs"
)

//	@title				Echo Web App
//	@version			1.0
//	@description	Echo based web application.

//	@contact.name		Sachin
//	@contact.url		https://www.sachinsharma.com/
//	@contact.email	trulysachin@gmail.com

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in													header
//	@name												Authorization

// @BasePath	/
func main() {
	cfg := config.NewConfig()

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	application.Start(cfg)
}
