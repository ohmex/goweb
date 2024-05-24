# GoWeb
A sample web MVC framework

## Get the echo swagger
go get -u github.com/swaggo/echo-swagger

## Install swag command
go install github.com/swaggo/swag/cmd/swag@latest

Run the Swag in your Go project root folder which contains main.go file, 
Swag will parse comments and generate required files(docs folder and docs/doc.go).

$ swag init

