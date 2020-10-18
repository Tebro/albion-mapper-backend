# Project moved

This project has been superceeded by: https://github.com/Portaler-Zone

# Albion mapper backend

Written in Go, using Go version 1.15


## Local setup

```
go mod download
docker-compose up -d db
source env.sh
go run main.go

# alternatively
docker-compose up -d
```


## Database migrations

The project uses the golang-migrate tool to handle migrations.

Install the CLI if you want to create new migrations: `brew install golang-migrate`

Then you can create migrations with: `migrate create -ext sql -dir migrations -seq <your_migration_name>`
