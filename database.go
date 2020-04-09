package main

import (
	"log"
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
			log.Println("Waiting for Database to be available #count", i, "/", maxTries)
			time.Sleep(1 * time.Second)
		}
	}

	return s
}

// DBCreateSchema creates database schema. Intended to be called by an admin command
func (s *Server) DBCreateSchema() *Server {

	// activate pgcrypto in order to use gen_random_uuid()
	if _, err := s.DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;"); err != nil {
		panic(err)
	}

	for _, model := range []interface{}{
		(*Super)(nil),
		(*Group)(nil),
		(*GroupSuper)(nil),
	} {
		log.Printf("Creating table for %T\n", model)
		err := s.DB.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}

	return s
}

// dbDropSchema should be used only by tests
func (s *Server) dbDropSchema() *Server {
	for _, model := range []interface{}{
		(*Super)(nil),
		(*Group)(nil),
		(*GroupSuper)(nil),
	} {
		err := s.DB.DropTable(model, &orm.DropTableOptions{IfExists: true})
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}
	return s
}

// DBMigrate performs database migrations (not implemented)
func (s *Server) DBMigrate() {
	log.Println("This will perform Database migrations")
}

// setupTestDatabase should be used only by tests
func (s *Server) setupEmptyTestDatabase() *Server {
	s.setupDatabase().dbDropSchema().DBCreateSchema() // first run with empty db
	return s
}
