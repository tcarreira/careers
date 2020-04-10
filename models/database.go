package models

import (
	"log"
	"os"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// SetupDatabase creates a DB connection and waits for availability. User is responsible for defer db.Close()
func SetupDatabase() *pg.DB {
	newDB := pg.Connect(&pg.Options{
		Addr:     getEnv("DB_HOST", "localhost") + ":" + getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASS", "password"),
		Database: getEnv("DB_NAME", "postgres"),
	})

	// wait for database to be ready
	maxTries := 30
	for i := 0; i < maxTries; i++ {
		_, err := newDB.Exec("SELECT 1")
		if err != nil {
			log.Println("Waiting for Database to be available #count", i, "/", maxTries)
			time.Sleep(1 * time.Second)
		}
	}

	return newDB
}

// SetupEmptyTestDatabase should be used only by tests
func SetupEmptyTestDatabase() *pg.DB {
	db := SetupDatabase()
	DropSchema(db)
	CreateSchema(db)

	return db
}

// CreateSchema creates database schema. Intended to be called by an admin command
func CreateSchema(db *pg.DB) {

	// activate pgcrypto in order to use gen_random_uuid()
	if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;"); err != nil {
		panic(err)
	}

	for _, model := range []interface{}{
		(*Super)(nil),
		(*Group)(nil),
		(*GroupSuper)(nil),
	} {
		log.Printf("Creating table for %T\n", model)
		err := db.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}
}

// DropSchema should be used only by tests
func DropSchema(db *pg.DB) {
	for _, model := range []interface{}{
		(*Super)(nil),
		(*Group)(nil),
		(*GroupSuper)(nil),
	} {
		err := db.DropTable(model, &orm.DropTableOptions{IfExists: true})
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}
}

// Migrate performs database migrations (not implemented)
func Migrate(db *pg.DB) {
	log.Println("This will perform Database migrations")
}
