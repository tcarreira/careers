// +build sql

package main

import (
	"testing"

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

func TestSuper_Create_badType(t *testing.T) {
	var err error

	super := Super{
		Type: "Something",
	}

	_, err = super.Create(s.DB)

	assert.Error(t, err)
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

func TestSuper_getByNameOrUUID_notExists(t *testing.T) {
	var got *Super
	var err error

	_, err = got.getByNameOrUUID(s.DB, "12356890")
	assert.Error(t, err)
}

func TestSuper_DeleteByNameOrUUID_byName(t *testing.T) {
	var super *Super
	var err error

	err = super.DeleteByNameOrUUID(s.DB, "Test1")
	assert.NoError(t, err)

	// Try again the same
	err = super.DeleteByNameOrUUID(s.DB, "Test1")
	assert.Error(t, err)
	assert.IsType(t, &errorSuperNotFound{""}, err)
}

func TestSuper_DeleteByNameOrUUID_byUUID(t *testing.T) {
	var super *Super
	var err error

	err = super.DeleteByNameOrUUID(s.DB, "47c0df01-a47d-497f-808d-181021f01c76")
	assert.NoError(t, err)

	// Try again the same
	err = super.DeleteByNameOrUUID(s.DB, "47c0df01-a47d-497f-808d-181021f01c76")
	assert.Error(t, err)
	assert.IsType(t, &errorSuperNotFound{""}, err)
}

func TestSuper_ReadAll_1(t *testing.T) {
	var got []Super
	var err error

	// Setup Database for this tests
	s.setupEmptyTestDatabase()

	supers := []Super{
		Super{
			Type: "HERO",
			Name: "h1",
			UUID: "47c0df01-a47d-497f-808d-181021f01c76",
		},
		Super{
			Type: "HERO",
			Name: "h2",
		},
		Super{
			Type: "HERO",
			Name: "h3",
		},
		Super{
			Type: "VILAN",
			Name: "v1",
		},
	}

	for i := range supers {
		supers[i].Create(s.DB)
	}

	// perform tests on previous data
	t.Run("Test ReadAll - no filters", func(t *testing.T) {
		super := Super{}
		got = super.ReadAll(s.DB)

		assert.NoError(t, err)
		assert.EqualValues(t, 4, len(got))

		expectedResuts := map[string]struct {
			Type string
			Name string
		}{
			"h1": {"HERO", "h1"},
			"h2": {"HERO", "h2"},
			"h3": {"HERO", "h3"},
			"v1": {"VILAN", "v1"},
		}

		for _, super := range got {
			assert.Equal(t, expectedResuts[super.Name].Type, super.Type)
			assert.Equal(t, expectedResuts[super.Name].Name, super.Name)
			delete(expectedResuts, super.Name)
		}
		assert.Equal(t, 0, len(expectedResuts))
	})

	t.Run("Test ReadAll - filter Type", func(t *testing.T) {
		super := Super{Type: "HERO"}
		got = super.ReadAll(s.DB)

		assert.NoError(t, err)
		assert.EqualValues(t, 3, len(got))

		expectedResuts := map[string]struct {
			Type string
			Name string
		}{
			"h1": {"HERO", "h1"},
			"h2": {"HERO", "h2"},
			"h3": {"HERO", "h3"},
		}

		for _, super := range got {
			assert.Equal(t, expectedResuts[super.Name].Type, super.Type)
			assert.Equal(t, expectedResuts[super.Name].Name, super.Name)
			delete(expectedResuts, super.Name)
		}
		assert.Equal(t, 0, len(expectedResuts))
	})

	t.Run("Test ReadAll - filter Name", func(t *testing.T) {
		super := Super{Name: "v1"}
		got = super.ReadAll(s.DB)

		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(got))

		expectedResuts := map[string]struct {
			Type string
			Name string
		}{
			"v1": {"VILAN", "v1"},
		}

		for _, super := range got {
			assert.Equal(t, expectedResuts[super.Name].Type, super.Type)
			assert.Equal(t, expectedResuts[super.Name].Name, super.Name)
			delete(expectedResuts, super.Name)
		}
		assert.Equal(t, 0, len(expectedResuts))
	})

	t.Run("Test ReadAll - filter UUID", func(t *testing.T) {
		super := Super{UUID: "47c0df01-a47d-497f-808d-181021f01c76"}
		got = super.ReadAll(s.DB)

		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(got))

		expectedResuts := map[string]struct {
			Type string
			Name string
			UUID string
		}{
			"h1": {"HERO", "h1", "47c0df01-a47d-497f-808d-181021f01c76"},
		}

		for _, super := range got {
			assert.Equal(t, expectedResuts[super.Name].Type, super.Type)
			assert.Equal(t, expectedResuts[super.Name].Name, super.Name)
			assert.Equal(t, expectedResuts[super.Name].UUID, super.UUID)
			delete(expectedResuts, super.Name)
		}
		assert.Equal(t, 0, len(expectedResuts))
	})

}
