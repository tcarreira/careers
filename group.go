package main

import "github.com/go-pg/pg"

// Group represents a group of supers
type Group struct {
	ID     uint64  `json:"-" sql:",pk"`
	Name   string  `json:"name" sql:",unique,notnull"`
	Supers []Super `json:"supers" pg:"many2many:group_supers"`
}

type errorGroupAlreadyExists struct {
	s string
}

func (e *errorGroupAlreadyExists) Error() string {
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

	return g, nil
}

// GetByName gets a group by its name
func (g *Group) GetByName(db *pg.DB, name string) (*Group, error) {
	return &Group{}, nil
}

// GetAllBySuper gets a list of Groups which Super is part of
func (g *Group) GetAllBySuper(db *pg.DB, super Super) ([]Group, error) {
	return []Group{}, nil
}
