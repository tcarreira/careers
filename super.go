package main

import (
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// SuperInterface interface for Super
type SuperInterface interface {
	Create(db *pg.DB) (*Super, error)
	Read(db *pg.DB) (*Super, error)
	ReadAll(db *pg.DB) ([]Super, error)
	Update(db *pg.DB) (*Super, error)
	Delete(db *pg.DB) error
}

// Super represents either a SuperHero or a SuperVilan
type Super struct {
	Type         string `json:"type" form:"type"`
	ID           uint64 `json:"id" sql:",pk"`
	Name         string `json:"name" form:"name" sql:",unique,notnull"`
	UUID         string `json:"uuid" form:"uuid" sql:",notnull,type:uuid,default:gen_random_uuid()"`
	FullName     string `json:"fullname"`
	Intelligence int64  `json:"intelligence"`
	Power        int64  `json:"power"`
	Occupation   string `json:"occupation"`
	ImageURL     string `json:"image_url"`
}

type errorSuperAlreadyExists struct {
	s string
}

func (e *errorSuperAlreadyExists) Error() string {
	return e.s
}

type errorSuperNotFound struct {
	s string
}

func (e *errorSuperNotFound) Error() string {
	return e.s
}

type errorSuperInvalidFields struct {
	s string
}

func (e *errorSuperInvalidFields) Error() string {
	return e.s
}

func (s *Super) validate() (*Super, error) {

	// Check Type is one of HERO|VILAN (should be enum...)
	switch strings.ToUpper(s.Type) {
	case
		"HERO",
		"VILAN":
		s.Type = strings.ToUpper(s.Type)
	default:
		return s, &errorSuperInvalidFields{"Type should be one of [\"HERO\", \"VILAN\"]"}
	}

	return s, nil
}

// Create saves the Super to database
func (s *Super) Create(db *pg.DB) (*Super, error) {
	if _, err := s.validate(); err != nil {
		return s, err
	}

	if err := db.Insert(s); err != nil {
		pgErr, ok := err.(pg.Error)
		if ok {
			if pgErr.IntegrityViolation() {
				return s, &errorSuperAlreadyExists{err.Error()}
			}
		}
		panic(err)
	}

	return s, nil
}

func (s *Super) getByNameOrUUID(db *pg.DB, idStr string) (*Super, error) {
	super := Super{}

	err := db.Model(&Super{}).
		Where("name = ?", idStr).
		WhereOr("upper(uuid::text) = ?", strings.ToUpper(idStr)).
		Select(&super)

	if err != nil {
		if err == pg.ErrNoRows {
			return &super, &errorSuperNotFound{err.Error()}
		}
		return &super, err
	}

	return &super, nil
}

// Read one Super from database
func (s *Super) Read(db *pg.DB) *Super {
	return s
}

// ReadAll read all Super from database (by ANDing super fields as filters)
func (s *Super) ReadAll(db *pg.DB) []Super {

	filter := func(q *orm.Query) (*orm.Query, error) {
		// Specs state filter only by Name and UUID
		if s.Type != "" {
			q = q.Where("upper(type) = ?", strings.ToUpper(superFilter.Type))
		}
		if s.Name != "" {
			// must match case
			q = q.Where("name = ?", superFilter.Name)
		}
		if s.UUID != "" {
			// postgres uuid is already case insensitive
			q = q.Where("upper(uuid::text) = ?", strings.ToUpper(superFilter.UUID))
		}

		// TODO: add other fields
	}

	// supersResult=[] instead of supersResult=nil
	supersResult := make([]Super, 0)

	err := db.Model(&supersResult).
		Apply(filter).
		Select()
	if err != nil {
		panic(err)
	}

	return supersResult

}

// Update a Super on database
func (s *Super) Update(db *pg.DB) *Super {
	return s
}

// Delete a super from database
func (s *Super) Delete(db *pg.DB) {
}
