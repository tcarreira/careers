// +build sql

package models

import (
	"testing"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseSimpleSelect(t *testing.T) {
	d := db.SetupDatabase()

	var num int

	// Simple params.
	_, err := s.DB.Query(pg.Scan(&num), "SELECT ?", 42)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, 42, num)
}

func TestDBCreateSchema(t *testing.T) {
	d := db.SetupDatabase()

	s.dbDropSchema()
	s.DBCreateSchema()
}
