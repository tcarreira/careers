package main

import (
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
	Type         string `json:"type"`
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	FullName     string `json:"fullname"`
	Intelligence string `json:"intelligence"`
	Power        string `json:"power"`
	Occupation   string `json:"occupation"`
	ImageURL     string `json:"image_url"`
}

type errorSuperAlreadyExists struct {
	s string
}

func (e *errorSuperAlreadyExists) Error() string {
	return e.s
}

// Create saves the Super to database
func (s *Super) Create(db *pg.DB) (*Super, error) {
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
