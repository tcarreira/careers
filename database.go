package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func (s *Server) setupDatabase() *Server {
	s.DB = pg.Connect(&pg.Options{
		Addr:     getEnv("DB_HOST", "localhost") + ":" + getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASS", "password"),
		Database: getEnv("DB_NAME", "postgres"),
	})

	// wait for database to be ready
	maxTries := 30
	for i := 0; i < maxTries; i++ {
		_, err := s.DB.Exec("SELECT 1")
		if err != nil {
			fmt.Println("Waiting for Database to be available #count", i, "/", maxTries)
			time.Sleep(1 * time.Second)
		}
	}

	return s
}

// DBCreateSchema creates database schema. Intended to be called by an admin command
func (s *Server) DBCreateSchema() {

	// activate pgcrypto in order to use gen_random_uuid()
	if _, err := s.DB.Exec("CREATE EXTENSION pgcrypto;"); err != nil {
		panic(err)
	}

	for _, model := range []interface{}{(*Super)(nil)} {
		fmt.Printf("Creating table for %T\n", model)
		err := s.DB.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

}

// DBMigrate performs database migrations (not implemented)
func (s *Server) DBMigrate() {
	fmt.Println("This will perform Database migrations")
}
