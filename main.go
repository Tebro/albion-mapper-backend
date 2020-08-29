package main

import (
	"log"
	"time"

	"github.com/tebro/albion-mapper-backend/albion"

	"github.com/tebro/albion-mapper-backend/server"

	"github.com/tebro/albion-mapper-backend/db"

	_ "github.com/go-sql-driver/mysql"
)

func setupDatabase() error {
	database, err := db.GetDb()
	if err != nil {
		return err
	}
	migrationsPath := "migrations"
	return db.RunMigrations(database, migrationsPath)
}

func main() {
	err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v\n", err)
	}

	go func() {
		log.Println("Starting cleaner")
		for {
			time.Sleep(time.Minute)
			log.Println("Cleaning up")
			err := albion.CleanupExpiredPortals()
			if err != nil {
				log.Printf("Error from cleanup: %v", err)
			}
		}
	}()

	err = server.StartServer()
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
