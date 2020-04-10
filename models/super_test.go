// +build sql

package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuper_Create(t *testing.T) {
	var got *Super
	var err error

	d := SetupEmptyTestDatabase()

	super1 := Super{
		Type: "Hero",
		Name: "Test1",
	}
	super2 := Super{
		Type:         "Vilan",
		Name:         "super2",
		UUID:         "47c0df01-a47d-497f-808d-181021f01c76",
		FullName:     "su per 2",
		Intelligence: 1,
		Power:        99,
		Occupation:   "something",
		ImageURL:     "url",
	}

	t.Run("TestSuper_Create - Create with Type and Name", func(t *testing.T) {

		got, err = super1.Create(d)

		assert.NoError(t, err)
		assert.Equal(t, "HERO", got.Type) // converts Hero -> HERO
		assert.Equal(t, "Test1", got.Name)
		assert.Less(t, uint64(0), got.ID) // got a new ID from database
		assert.NotEmpty(t, got.UUID)      // got new UUID from database
	})

	t.Run("TestSuper_Create - Repeat create - should fail", func(t *testing.T) {
		// Try again the same
		_, err = super1.Create(d)
		assert.Error(t, err)
	})

	t.Run("TestSuper_Create - Create with all fields", func(t *testing.T) {

		got, err = super2.Create(d)

		assert.NoError(t, err)
		assert.Equal(t, "VILAN", got.Type) // converts Hero -> HERO
		assert.Equal(t, "super2", got.Name)
		assert.Less(t, uint64(0), got.ID)                                 // got a new id from database
		assert.Equal(t, "47c0df01-a47d-497f-808d-181021f01c76", got.UUID) // got new UUID from database
		assert.Equal(t, "su per 2", got.FullName)
		assert.EqualValues(t, 1, got.Intelligence)
		assert.EqualValues(t, 99, got.Power)
		assert.Equal(t, "something", got.Occupation)
		assert.Equal(t, "url", got.ImageURL)
	})

	t.Run("TestSuper_Create - try a bad Super.Type", func(t *testing.T) {

		badSuper := Super{
			Type: "Something",
		}

		_, err = badSuper.Create(d)

		assert.Error(t, err)
	})
}

func TestSuper_GetByNameOrUUID_byUUID(t *testing.T) {
	var got *Super
	var err error

	d := SetupEmptyTestDatabase()

	supers := []Super{
		{
			Type: "Hero",
			Name: "Test1",
		},
		{
			Type: "Vilan",
			Name: "super2",
			UUID: "40000001-a47d-497f-808d-181021f01c76",
		},
	}

	for i := range supers {
		supers[i].Create(d)
	}

	t.Run("TestSuper_GetByNameOrUUID - by UUID", func(t *testing.T) {
		got, err = got.GetByNameOrUUID(d, "40000001-a47d-497f-808d-181021f01c76")
		assert.NoError(t, err)
		assert.Equal(t, "VILAN", got.Type)
		assert.Equal(t, "super2", got.Name)
	})

	t.Run("TestSuper_GetByNameOrUUID - by Name", func(t *testing.T) {
		got, err = got.GetByNameOrUUID(d, "Test1")
		assert.NoError(t, err)
		assert.Equal(t, "HERO", got.Type)
		assert.Equal(t, "Test1", got.Name)
	})

	t.Run("TestSuper_GetByNameOrUUID - does not exist", func(t *testing.T) {
		_, err = got.GetByNameOrUUID(d, "12356890")
		assert.Error(t, err)
	})
}

func TestSuper_DeleteByNameOrUUID(t *testing.T) {
	var err error

	d := SetupEmptyTestDatabase()

	supers := []Super{
		{
			Type: "Hero",
			Name: "Test1",
		},
		{
			Type: "Vilan",
			Name: "super2",
			UUID: "40000002-a47d-497f-808d-181021f01c76",
		},
	}

	for i := range supers {
		supers[i].Create(d)
	}

	t.Run("TestSuper_DeleteByNameOrUUID - by Name", func(t *testing.T) {
		err = new(Super).DeleteByNameOrUUID(d, "Test1")
		assert.NoError(t, err)
	})

	t.Run("TestSuper_DeleteByNameOrUUID - by Name - same again", func(t *testing.T) {
		err = new(Super).DeleteByNameOrUUID(d, "Test1")
		assert.Error(t, err)
		assert.IsType(t, &ErrorSuperNotFound{""}, err)
	})

	t.Run("TestSuper_DeleteByNameOrUUID - by UUID", func(t *testing.T) {
		err = new(Super).DeleteByNameOrUUID(d, "40000002-a47d-497f-808d-181021f01c76")
		assert.NoError(t, err)
	})
	t.Run("TestSuper_DeleteByNameOrUUID - by UUID - same again", func(t *testing.T) {
		// Try again the same
		err = new(Super).DeleteByNameOrUUID(d, "40000002-a47d-497f-808d-181021f01c76")
		assert.Error(t, err)
		assert.IsType(t, &ErrorSuperNotFound{""}, err)
	})
}

func TestSuper_ReadAll(t *testing.T) {
	var got []Super
	var err error

	d := SetupEmptyTestDatabase()

	supers := []Super{
		{
			Type: "HERO",
			Name: "h1",
			UUID: "47c0df01-a47d-497f-808d-181021f01c76",
		},
		{
			Type: "HERO",
			Name: "h2",
		},
		{
			Type: "HERO",
			Name: "h3",
		},
		{
			Type: "VILAN",
			Name: "v1",
		},
	}

	for i := range supers {
		supers[i].Create(d)
	}

	// perform tests on previous data
	t.Run("Test ReadAll - no filters", func(t *testing.T) {
		super := Super{}
		got = super.ReadAll(d)

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
		got = super.ReadAll(d)

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
		got = super.ReadAll(d)

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
		got = super.ReadAll(d)

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

func TestSuper_GroupsRelatives(t *testing.T) {
	d := SetupEmptyTestDatabase()

	supers := []Super{
		{Type: "HERO", Name: "main", UUID: "41f0bc0e-89f7-4ea7-a4f5-9d08e5383b9c"},
		{Type: "HERO", Name: "rel1"},
		{Type: "HERO", Name: "rel2"},
		{Type: "HERO", Name: "rel3"},
	}
	for i := range supers {
		supers[i].Create(d)
	}
	groups := []Group{
		{Name: "g1", Supers: supers[0:3]},
		{Name: "g2", Supers: supers[1:4]},
		{Name: "g3", Supers: supers[0:4]},
	}
	for i := range groups {
		groups[i].Create(d)
	}

	t.Run("TestSuper_GroupsRelatives - Marshal Super JSON", func(t *testing.T) {
		// This test is very prone to errors
		// special attention to "relatives_count": and "groups":
		super, _ := new(Super).GetByNameOrUUID(d, supers[0].Name)

		var superJSON []byte
		superJSON, err := json.Marshal(super)

		assert.NoError(t, err)
		assert.Equal(t,
			string(`{"uuid":"41f0bc0e-89f7-4ea7-a4f5-9d08e5383b9c","type":"HERO","name":"main","fullname":"","intelligence":"0","power":"0","occupation":"","image_url":"","groups":["g1","g3"],"relatives_count":"3"}`),
			string(superJSON),
		)

	})
	t.Run("TestSuper_GroupsRelatives - RelativesCount", func(t *testing.T) {
		super, err := new(Super).GetByNameOrUUID(d, supers[0].Name)

		assert.NoError(t, err)
		assert.Equal(t, int(3), super.RelativesCount)

	})

	t.Run("TestSuper_GroupsRelatives - Marshal Super JSON", func(t *testing.T) {
		// This test is very prone to errors
		// special attention to "relatives_count": and "groups":
		super := &Super{Name: supers[0].Name}
		superList := super.ReadAll(d)

		var superListJSON []byte
		superListJSON, err := json.Marshal(superList)

		assert.NoError(t, err)
		assert.Equal(t,
			string(`[{"uuid":"41f0bc0e-89f7-4ea7-a4f5-9d08e5383b9c","type":"HERO","name":"main","fullname":"","intelligence":"0","power":"0","occupation":"","image_url":"","groups":["g1","g3"],"relatives_count":"3"}]`),
			string(superListJSON),
		)

	})

}
