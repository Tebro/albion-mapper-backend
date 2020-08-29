package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"

	// Migration specific dependency
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB

// GetDb returns the database connection
func GetDb() (*sql.DB, error) {
	if db == nil {
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		host := os.Getenv("MYSQL_HOST")
		port := os.Getenv("MYSQL_PORT")
		dbName := os.Getenv("MYSQL_DATABASE")

		url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", user, password, host, port, dbName)
		res, err := sql.Open("mysql", url)
		if err != nil {
			return nil, err
		}
		db = res
	}
	return db, nil
}

// RunMigrations runs the Database migrations
func RunMigrations(db *sql.DB, migrationsPath string) error {
	log.Println("Running database migrations")
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		os.Getenv("MYSQL_DATABASE"),
		driver)

	if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		// Having no changes is actually an error with Up, but not for us.
		log.Println("Database already up to date")
		return nil
	}
	log.Println("Database migration completed")
	return err
}

// Hello is a simple testing function
func Hello() (*sql.Rows, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT 'hello';")

	return rows, err
}
