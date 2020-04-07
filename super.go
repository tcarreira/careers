package main

import "github.com/go-pg/pg"

// SuperInterface interface for Super
type SuperInterface interface {
	Create(db *pg.DB) *Super
	Read(db *pg.DB) *Super
	ReadAll(db *pg.DB) []Super
	Update(db *pg.DB) *Super
	Delete(db *pg.DB)
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

// Create saves the Super to database
func (s *Super) Create(db *pg.DB) *Super {
	return s
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
