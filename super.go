package main

import (
	"encoding/json"
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
	ID           uint64  `json:"-" sql:",pk"`
	UUID         string  `json:"uuid" form:"uuid" sql:",notnull,type:uuid,default:gen_random_uuid()"`
	Type         string  `json:"type" form:"type"`
	Name         string  `json:"name" form:"name" sql:",unique,notnull"`
	FullName     string  `json:"fullname"`
	Intelligence int64   `json:"intelligence"`
	Power        int64   `json:"power"`
	Occupation   string  `json:"occupation"`
	ImageURL     string  `json:"image_url"`
	Groups       []Group `json:"-" pg:"many2many:group_supers,joinFK:group_id"`
}

// MarshalJSON will render a Super JSON with a []string of group names instead of []Group
func (s *Super) MarshalJSON() ([]byte, error) {
	type Alias Super

	groupsNames := make([]string, 0)
	for _, group := range s.Groups {
		groupsNames = append(groupsNames, group.Name)
	}

	return json.Marshal(&struct {
		*Alias
		Groups []string `json:"groups"`
	}{
		Alias:  (*Alias)(s),
		Groups: groupsNames,
	})
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

	err := db.Model(&super).
		Relation("Groups").
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
			q = q.Where("upper(type) = ?", strings.ToUpper(s.Type))
		}
		if s.Name != "" {
			// must match case
			q = q.Where("name = ?", s.Name)
		}
		if s.UUID != "" {
			// postgres uuid is already case insensitive
			q = q.Where("upper(uuid::text) = ?", strings.ToUpper(s.UUID))
		}

		// TODO: add other fields

		return q, nil
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

// DeleteByNameOrUUID deletes Super from database, using name or uuid
func (s *Super) DeleteByNameOrUUID(db *pg.DB, idStr string) error {

	res, err := db.Model(&Super{}).
		Where("name = ?", idStr).
		WhereOr("upper(uuid::text) = ?", strings.ToUpper(idStr)).
		Delete()

	if err != nil {
		if err == pg.ErrNoRows {
			return &errorSuperNotFound{err.Error()}
		}
		return err
	}
	if res.RowsAffected() < 1 {
		return &errorSuperNotFound{"Can't delete Super - Not Found"}
	}

	return nil
}

// Delete a super from database
func (s *Super) Delete(db *pg.DB) {

}
