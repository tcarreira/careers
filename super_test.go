// +build sql

package main

import (
	"testing"

	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
)

var s *Server

func TestSuper_Create_1(t *testing.T) {
	var got *Super
	var err error

	// Setup DB on first test
	s = &Server{}
	s.setupEmptyTestDatabase()

	super := Super{
		Type: "Hero",
		Name: "Test1",
	}

	got, err = super.Create(s.DB)

	assert.NoError(t, err)
	assert.Equal(t, "HERO", got.Type) // converts Hero -> HERO
	assert.Equal(t, "Test1", got.Name)
	assert.NotEqual(t, 0, got.ID) // got a new ID from database
	assert.NotEmpty(t, got.UUID)  // got new UUID from database

	// Try again the same
	got, err = super.Create(s.DB)
	assert.Error(t, err)
}

func TestSuper_Create_1repeated(t *testing.T) {
	var err error

	super := Super{
		Type: "Hero",
		Name: "Test1",
	}

	// Try again the same
	_, err = super.Create(s.DB)
	assert.Error(t, err)
}

func TestSuper_Create_2(t *testing.T) {
	var got *Super
	var err error

	super := Super{
		Type:         "Vilan",
		Name:         "super2",
		UUID:         "47c0df01-a47d-497f-808d-181021f01c76",
		FullName:     "su per 2",
		Intelligence: 1,
		Power:        99,
		Occupation:   "something",
		ImageURL:     "url",
	}

	got, err = super.Create(s.DB)

	assert.NoError(t, err)
	assert.Equal(t, "VILAN", got.Type) // converts Hero -> HERO
	assert.Equal(t, "super2", got.Name)
	assert.Less(t, uint64(0), got.ID)                                 // got a new ID from database
	assert.Equal(t, "47c0df01-a47d-497f-808d-181021f01c76", got.UUID) // got new UUID from database
	assert.Equal(t, "su per 2", got.FullName)
	assert.EqualValues(t, 1, got.Intelligence)
	assert.EqualValues(t, 99, got.Power)
	assert.Equal(t, "something", got.Occupation)
	assert.Equal(t, "url", got.ImageURL)

}

func TestSuper_getByNameOrUUID_byUUID(t *testing.T) {
	var got *Super
	var err error

	got, err = got.getByNameOrUUID(s.DB, "47c0df01-a47d-497f-808d-181021f01c76")
	assert.NoError(t, err)
	assert.Equal(t, "VILAN", got.Type)
	assert.Equal(t, "super2", got.Name)
}

func TestSuper_getByNameOrUUID_byName(t *testing.T) {
	var got *Super
	var err error

	got, err = got.getByNameOrUUID(s.DB, "Test1")
	assert.NoError(t, err)
	assert.Equal(t, "HERO", got.Type)
	assert.Equal(t, "Test1", got.Name)
}

func TestSuper_Delete(t *testing.T) {
	type fields struct {
		Type         string
		ID           uint64
		Name         string
		UUID         string
		FullName     string
		Intelligence int64
		Power        int64
		Occupation   string
		ImageURL     string
	}
	type args struct {
		db *pg.DB
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Super{
				Type:         tt.fields.Type,
				ID:           tt.fields.ID,
				Name:         tt.fields.Name,
				UUID:         tt.fields.UUID,
				FullName:     tt.fields.FullName,
				Intelligence: tt.fields.Intelligence,
				Power:        tt.fields.Power,
				Occupation:   tt.fields.Occupation,
				ImageURL:     tt.fields.ImageURL,
			}
			s.Delete(tt.args.db)
		})
	}
}
