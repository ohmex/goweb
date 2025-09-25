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

This project supports PostgreSQL table partitioning by the `domain` field for multi-tenancy with **automatic partition creation** for each domain.

## Quick Start

1. **Enable Partitioning**: Set the environment variable:
   ```bash
   export DB_PARTITIONING_ENABLED=true
   ```

2. **Create Domains**: Use the API to create domains - partitions are created automatically:
   ```bash
   curl -X POST http://localhost:8080/api/domain \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{"name": "NewCompany"}'
   ```

3. **Run Migrations**: Partitions are created automatically during migration:
   ```bash
   go run main.go migrate
   ```

## Features

- ✅ **Automatic Partition Creation**: Partitions created automatically when domains are created
- ✅ **Multi-tenant Isolation**: Complete data separation between domains
- ✅ **Performance Optimization**: Improved query performance for domain-specific data
- ✅ **Migration Integration**: Works with existing migration system
- ✅ **Error Handling**: Graceful handling of partitioning failures

## Detailed Documentation

For comprehensive information about the partitioning system, see [docs/partitioning.md](docs/partitioning.md).

## Configuration

- **Environment Variable**: `DB_PARTITIONING_ENABLED=true`
- **Database**: PostgreSQL or YugabyteDB (`DB_DRIVER=postgres` or `DB_DRIVER=yugabytedb`)
- **Affected Tables**: Tables using `BaseResource` struct (currently `posts`)
- **Partition Type**: LIST partitions by domain UUID
- **Partitioning Support**: Both PostgreSQL and YugabyteDB support table partitioning

## Supported Databases

This application supports multiple database backends:

- **MySQL**: `DB_DRIVER=mysql`
- **PostgreSQL**: `DB_DRIVER=postgres`
- **YugabyteDB**: `DB_DRIVER=yugabytedb` (PostgreSQL-compatible distributed SQL database)

### YugabyteDB Setup

YugabyteDB is a distributed SQL database that is PostgreSQL-compatible. To use YugabyteDB:

1. **Set Environment Variables**:
   ```bash
   export DB_DRIVER=yugabytedb
   export DB_HOST=localhost
   export DB_PORT=5433
   export DB_USER=yugabyte
   export DB_PASSWORD=yugabyte
   export DB_NAME=yugabyte
   export DB_PARTITIONING_ENABLED=true  # Enable partitioning for multi-tenancy
   ```

2. **Start YugabyteDB with Docker Compose**:
   ```bash
   docker-compose up echo_yugabytedb
   ```

3. **Access YugabyteDB**:
   - YSQL (PostgreSQL-compatible): `localhost:5433`
   - Master Web UI: `localhost:7000`
   - TServer Web UI: `localhost:9000`

4. **Partitioning Support**:
   - YugabyteDB supports PostgreSQL-compatible table partitioning
   - Automatic partition creation for each domain
   - LIST partitioning by domain UUID for multi-tenant isolation

