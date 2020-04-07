package main

import (
	"strings"

	"github.com/go-pg/pg"
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

// Read queries one Super from database
func (s *Super) Read(db *pg.DB) *Super {
	return s
}

// ReadAll read all Super from database
func (s *Super) ReadAll(db *pg.DB) []Super {
	return []Super{}
}

// Update a Super on database
func (s *Super) Update(db *pg.DB) *Super {
	return s
}

// Delete a super from database
func (s *Super) Delete(db *pg.DB) {
}
