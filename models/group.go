package models

import (
	"encoding/json"
	"strings"

	"github.com/go-pg/pg"
)

// Group represents a group of supers
type Group struct {
	tableName  struct{} `pg:"alias:g"`
	ID         uint64   `json:"-" sql:",pk"`
	Name       string   `json:"name" sql:",unique,notnull"`
	Supers     []Super  `json:"-" pg:"many2many:group_supers,joinFK:super_id"`
	SupersList []string `json:"supers" sql:"-" `
}

// UnmarshalJSON will instantiate a Group from a JSON, where Supers is a []string of Super names
func (g *Group) UnmarshalJSON(data []byte) error {
	type Alias Group
	aux := &struct {
		*Alias
		Supers []string `json:"supers"`
	}{
		Alias: (*Alias)(g),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	supers := make([]Super, 0)
	for _, name := range aux.Supers {
		supers = append(supers, Super{Name: name})
	}
	g.Supers = supers

	return nil
}

// GroupSuper represents many2many table Groups-Supers
type GroupSuper struct {
	tableName struct{} `pg:"alias:g2s"`
	GroupID   uint64   `sql:"group_id,pk"`
	Group     *Group
	SuperID   uint64 `sql:"super_id,pk"`
	Super     *Super
}

// ErrorGroupAlreadyExists Group Already Exists - extends error
type ErrorGroupAlreadyExists struct {
	s string
}

func (e *ErrorGroupAlreadyExists) Error() string {
	return e.s
}

// ErrorGroupSuperRelation  Group Super Relation - extends error
type ErrorGroupSuperRelation struct {
	s string
}

func (e *ErrorGroupSuperRelation) Error() string {
	return e.s
}

// ErrorGroupNotFound Group Not Found - extends error
type ErrorGroupNotFound struct {
	s string
}

func (e *ErrorGroupNotFound) Error() string {
	return e.s
}

// Create a group with a list of Supers
func (g *Group) Create(db *pg.DB) (*Group, error) {
	if err := db.Insert(g); err != nil {
		pgErr, ok := err.(pg.Error)
		if ok {
			if pgErr.IntegrityViolation() {
				return g, &ErrorGroupAlreadyExists{err.Error()}
			}
		}
		panic(err)
	}

	var minorErrors []string
	for _, super := range g.Supers {
		var s *Super
		var err error

		// Get Super from DB
		s, err = super.GetByNameOrUUID(db, super.Name)
		if err != nil {
			minorErrors = append(minorErrors,
				"(super:'"+super.Name+"') "+err.Error(),
			)
		} else {
			if err = db.Insert(&GroupSuper{
				GroupID: g.ID,
				Group:   g,
				SuperID: s.ID,
				Super:   s,
			}); err != nil {
				minorErrors = append(minorErrors, err.Error())
			}
		}
	}

	if len(minorErrors) > 0 {
		return g, &ErrorGroupSuperRelation{strings.Join(minorErrors, " | ")}
	}

	return g, nil
}

// GetByName gets a group by its name
func (g *Group) GetByName(db *pg.DB, name string) (*Group, error) {
	group := Group{}

	err := db.Model(&group).
		Relation("Supers").
		Where("name = ?", name).
		Select(&group)

	if err != nil {
		if err == pg.ErrNoRows {
			return &group, &ErrorGroupNotFound{err.Error()}
		}
		return &group, err
	}

	// create the Super Names List as []string
	for _, super := range group.Supers {
		group.SupersList = append(group.SupersList, super.Name)
	}

	return &group, nil
}

// GetAllBySuper gets a list of Groups which Super is part of
func (g *Group) GetAllBySuper(db *pg.DB, super Super) ([]Group, error) {
	var results []Group

	err := db.Model(&results).
		Join("JOIN group_supers AS gs").
		JoinOn("gs.group_id = \"group\".id").
		JoinOn("gs.super_id = ?", super.ID).
		Relation("Supers").
		Select()

	if err != nil {
		return results, err
	}

	return results, nil
}
