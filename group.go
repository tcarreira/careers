package main

import (
	"strings"

	"github.com/go-pg/pg"
)

// Group represents a group of supers
type Group struct {
	tableName struct{} `pg:"alias:g"`
	ID        uint64   `json:"-" sql:",pk"`
	Name      string   `json:"name" sql:",unique,notnull"`
	Supers    []Super  `json:"supers" pg:"many2many:group_supers,joinFK:super_id"`
}

// GroupSuper represents many2many table Groups-Supers
type GroupSuper struct {
	tableName struct{} `pg:"alias:g2s"`
	GroupID   uint64   `sql:"group_id,pk"`
	Group     *Group
	SuperID   uint64 `sql:"super_id,pk"`
	Super     *Super
}

type errorGroupAlreadyExists struct {
	s string
}

func (e *errorGroupAlreadyExists) Error() string {
	return e.s
}

type errorGroupSuperRelation struct {
	s string
}

func (e *errorGroupSuperRelation) Error() string {
	return e.s
}

type errorGroupNotFound struct {
	s string
}

func (e *errorGroupNotFound) Error() string {
	return e.s
}

// Create a group with a list of Supers
func (g *Group) Create(db *pg.DB) (*Group, error) {
	if err := db.Insert(g); err != nil {
		pgErr, ok := err.(pg.Error)
		if ok {
			if pgErr.IntegrityViolation() {
				return g, &errorGroupAlreadyExists{err.Error()}
			}
		}
		panic(err)
	}

	var minorErrors []string
	for _, super := range g.Supers {
		if err := db.Insert(&GroupSuper{
			GroupID: g.ID,
			Group:   g,
			SuperID: super.ID,
			Super:   &super,
		}); err != nil {
			minorErrors = append(minorErrors, err.Error())
		}
	}

	if len(minorErrors) > 0 {
		return g, &errorGroupSuperRelation{strings.Join(minorErrors, " | ")}
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
			return &group, &errorGroupNotFound{err.Error()}
		}
		return &group, err
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
		if err == pg.ErrNoRows {
			return results, &errorGroupNotFound{err.Error()}
		}
		return results, err
	}

	return results, nil
}
