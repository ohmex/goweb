# GoWeb
A sample web MVC framework

## Get the echo swagger
go get -u github.com/swaggo/echo-swagger

## Install swag command
go install github.com/swaggo/swag/cmd/swag@latest

Run the Swag in your Go project root folder which contains main.go file, 
Swag will parse comments and generate required files(docs folder and docs/doc.go).

$ swag init

# Goal - to build a web framework that supports:
1. Multi tenant ecosystem
2. REST API based resource management
3. RBAC access controls based on Casbin
4. JWT based authentication with custom claims

## TODO:
1. Token should be unique for the User & Domain, No need to send domain explicitly, it can be the part of token

