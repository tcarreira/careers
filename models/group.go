package models

import (
	"encoding/json"
	"strings"

	"github.com/go-pg/pg/v9"
)

// Group represents a group of supers
type Group struct {
	tableName  struct{} `json:"-" pg:"superhero_groups,alias:g"` // json tag for swaggo bug
	ID         uint64   `json:"-" pg:",pk"`
	Name       string   `json:"name" example:"group1" pg:",unique,notnull"`
	Supers     []Super  `json:"-" pg:"many2many:superhero_group_supers,joinFK:super_id"`
	SupersList []string `json:"supers,nilasempty" example:"name1" pg:"-" `
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
	tableName struct{} `pg:"superhero_group_supers,alias:g2s"`
	GroupID   uint64   `pg:"group_id,pk"`
	Group     *Group
	SuperID   uint64 `pg:"super_id,pk"`
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

	g.SupersList = make([]string, 0) // empty array instead of null

	var minorErrors []string
	for i := range g.Supers {
		var err error
		s := &(g.Supers[i])

		// Get Super from DB
		superName := s.Name // we need to keep the name in case of entity not found
		s, err = s.GetByNameOrUUID(db, superName)
		if err != nil {
			minorErrors = append(minorErrors,
				"(super:'"+superName+"') "+err.Error(),
			)
		} else {
			if err = db.Insert(&GroupSuper{
				GroupID: g.ID,
				Group:   g,
				SuperID: s.ID,
				Super:   s,
			}); err != nil {
				minorErrors = append(minorErrors, err.Error())
			} else {
				// Really commited the transaction
				g.SupersList = append(g.SupersList, s.Name)
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
		Join("JOIN superhero_group_supers AS gs").
		JoinOn("gs.group_id = g.id").
		JoinOn("gs.super_id = ?", super.ID).
		Relation("Supers").
		Select()

	if err != nil {
		return results, err
	}

	return results, nil
}
