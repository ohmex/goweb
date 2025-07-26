# GoWeb
A sample web MVC framework

## Get the echo swagger
go get -u github.com/swaggo/echo-swagger

## Insomnia
(metadata, body) => JSON.parse(atob(JSON.parse(body).accessToken.split(`.`)[1])).domains[0].UUID

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

# Database Partitioning Support

This project supports PostgreSQL table partitioning by the `domain` field for multi-tenancy. Partitioning is controlled by the `DB_PARTITIONING_ENABLED` environment variable.

## How to Enable Partitioning

- Set the following environment variable before running migrations or starting the application:

```
DB_PARTITIONING_ENABLED=true
```

- Partitioning is only applied if the database driver is PostgreSQL (`DB_DRIVER=postgres`).
- When enabled, the `posts` table will be partitioned by the `domain` column during migration.

## Effect
- Improves performance and scalability for multi-tenant data separation.
- Only affects tables that use the `BaseResource` struct (currently, only `posts`).
- If disabled or using a non-PostgreSQL database, standard tables are created without partitioning.

