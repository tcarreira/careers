// +build sql

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertSupers(t *testing.T, expected, actual []Super) {

	type testedSuper struct {
		Type string
		Name string
	}

	expectedResuts := make(map[string]testedSuper, len(expected))
	for _, val := range expected {
		expectedResuts[val.Name] = testedSuper{val.Type, val.Name}
	}

	for _, val := range actual {
		assert.Equal(t, expectedResuts[val.Name].Type, val.Type)
		assert.Equal(t, expectedResuts[val.Name].Name, val.Name)
		delete(expectedResuts, val.Name)
	}

	assert.Equal(t, 0, len(expectedResuts))

}

func assertGroups(t *testing.T, expected, actual []Group) {

	expectedResuts := make(map[string]Group, len(expected))
	for _, val := range expected {
		expectedResuts[val.Name] = val
	}

	for _, val := range actual {
		assert.Equal(t, expectedResuts[val.Name].Name, val.Name)
		assertSupers(t, expectedResuts[val.Name].Supers, val.Supers)
		delete(expectedResuts, val.Name)
	}

	assert.Equal(t, 0, len(expectedResuts))

}

func TestGroup_Create(t *testing.T) {
	var got *Group
	var err error

	d := SetupEmptyTestDatabase()

	supers := []Super{
		{Type: "HERO", Name: "t1"},
		{Type: "HERO", Name: "t2"},
	}
	for i := range supers {
		supers[i].Create(d)
	}

	group := Group{
		Name:   "group1",
		Supers: supers,
	}

	t.Run("TestGroup_Create - create a group", func(t *testing.T) {
		got, err = group.Create(d)

		assert.NoError(t, err)
		assert.Equal(t, "group1", got.Name)
		assertSupers(t, supers, got.Supers)
	})

	t.Run("TestGroup_Create - create a group - same again", func(t *testing.T) {
		got, err = group.Create(d)
		assert.Error(t, err)
	})

}

func TestGroup_GetByName(t *testing.T) {

	d := SetupEmptyTestDatabase()

	supers := []Super{
		{Type: "HERO", Name: "t1"},
		{Type: "HERO", Name: "t2"},
	}
	for i := range supers {
		supers[i].Create(d)
	}

	groups := []Group{
		{
			Name:   "group1",
			Supers: supers,
		},
		{
			Name:   "group2",
			Supers: []Super{supers[0]},
		},
		{
			Name:   "group3",
			Supers: []Super{},
		},
	}
	for i := range groups {
		groups[i].Create(d)
	}

	t.Run("TestGroup_GetByName - test1 - with multiple Super", func(t *testing.T) {
		got, err := new(Group).GetByName(d, "group1")

		assert.NoError(t, err)
		assert.Equal(t, "group1", got.Name)
		assertSupers(t, supers, got.Supers)
	})

	t.Run("TestGroup_GetByName - test2 - with one Super", func(t *testing.T) {
		got, err := new(Group).GetByName(d, "group2")

		assert.NoError(t, err)
		assert.Equal(t, "group2", got.Name)
		assertSupers(t, []Super{{Type: "HERO", Name: "t1"}}, got.Supers)
	})

	t.Run("TestGroup_GetByName - test3 - with zero Super", func(t *testing.T) {
		got, err := new(Group).GetByName(d, "group3")

		assert.NoError(t, err)
		assert.Equal(t, "group3", got.Name)
		assertSupers(t, []Super{}, got.Supers)
	})

	t.Run("TestGroup_GetByName - not found", func(t *testing.T) {
		_, err := new(Group).GetByName(d, "groupX")
		assert.Error(t, err)
	})
}

func TestGroup_GetAllBySuper(t *testing.T) {
	d := SetupEmptyTestDatabase()

	supers := []Super{
		{Type: "HERO", Name: "s1"},
		{Type: "HERO", Name: "s2"},
		{Type: "HERO", Name: "s3"},
	}
	for i := range supers {
		supers[i].Create(d)
	}

	groups := []Group{
		{
			Name:   "group1",
			Supers: supers[0:2],
		},
		{
			Name:   "group2",
			Supers: []Super{supers[0]},
		},
		{
			Name:   "group3",
			Supers: []Super{},
		},
	}
	for i := range groups {
		groups[i].Create(d)
	}

	t.Run("TestGroup_GetAllBySuper - multiple groups", func(t *testing.T) {
		got, err := new(Group).GetAllBySuper(d, supers[0])

		assert.NoError(t, err)
		assert.Equal(t, 2, len(got))
		assertGroups(t, groups[0:2], got)

	})

	t.Run("TestGroup_GetAllBySuper - one group", func(t *testing.T) {
		got, err := new(Group).GetAllBySuper(d, supers[1])

		assert.NoError(t, err)
		assert.Equal(t, 1, len(got))
		assertGroups(t, []Group{groups[0]}, got)
	})

	t.Run("TestGroup_GetAllBySuper - zero groups", func(t *testing.T) {
		got, err := new(Group).GetAllBySuper(d, supers[2])

		assert.NoError(t, err)
		assert.Equal(t, 0, len(got))
	})

}
