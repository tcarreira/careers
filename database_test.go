// +build sql

package main

import (
	"testing"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseSimpleSelect(t *testing.T) {
	s := Server{}
	s.setupDatabase()

	var num int

	// Simple params.
	_, err := s.DB.Query(pg.Scan(&num), "SELECT ?", 42)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, 42, num)
}

func TestDBCreateSchema(t *testing.T) {
	s := new(Server)
	s.setupDatabase()

	s.dbDropSchema()
	s.DBCreateSchema()
}
