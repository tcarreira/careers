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
// swagger:model Super
type Super struct {
	ID             uint64   `json:"-" sql:",pk"`
	UUID           string   `json:"uuid" example:"47c0df01-a47d-497f-808d-181021f01c76" form:"uuid" sql:",notnull,type:uuid,default:gen_random_uuid()"`
	Type           string   `json:"type" example:"HERO" enums:"HERO,VILAN" form:"type"`
	Name           string   `json:"name" form:"name" example:"SuperHero1" sql:",unique,notnull"`
	FullName       string   `json:"fullname" example:"SuperHero1's Full Name"`
	Intelligence   int64    `json:"intelligence,string" example:"90"`
	Power          int64    `json:"power,string" example:"80"`
	Occupation     string   `json:"occupation" example:"Programmer"`
	ImageURL       string   `json:"image_url" example:"https://http.cat/200"`
	Groups         []Group  `json:"-" pg:"many2many:group_supers,joinFK:group_id"`
	GroupsList     []string `json:"groups" example:"group1,group2" sql:"-"`
	RelativesCount int      `json:"relatives_count,string" sql:"-"`
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

	err := db.Model(&super).TableExpr("supers AS s").
		Relation("Groups").
		Column("s.*").ColumnExpr("count(distinct relatives.id) AS relatives_count").
		Join("LEFT JOIN group_supers AS s2g ON s.id = s2g.super_id").
		Join("LEFT JOIN group_supers AS g2s ON s2g.group_id = g2s.group_id").
		Join("LEFT JOIN supers AS relatives ON g2s.super_id = relatives.id AND g2s.super_id != s.id").
		Where("s.name = ?", idStr).
		WhereOr("upper(s.uuid::text) = ?", strings.ToUpper(idStr)).
		Group("s.id").
		Select(&super)

	if err != nil {
		if err == pg.ErrNoRows {
			return &super, &errorSuperNotFound{err.Error()}
		}
		return &super, err
	}

	// create the Group Names List as []string
	for _, group := range super.Groups {
		super.GroupsList = append(super.GroupsList, group.Name)
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
			q = q.Where("upper(s.type) = ?", strings.ToUpper(s.Type))
		}
		if s.Name != "" {
			// must match case
			q = q.Where("s.name = ?", s.Name)
		}
		if s.UUID != "" {
			// postgres uuid is already case insensitive
			q = q.Where("upper(s.uuid::text) = ?", strings.ToUpper(s.UUID))
		}

		// TODO: add other fields

		return q, nil
	}

	// supersResult=[] instead of supersResult=nil
	supersResult := make([]Super, 0)

	err := db.Model(&supersResult).TableExpr("supers AS s").
		Relation("Groups").
		Column("s.*").ColumnExpr("count(distinct relatives.id) AS relatives_count").
		Join("LEFT JOIN group_supers AS s2g ON s.id = s2g.super_id").
		Join("LEFT JOIN group_supers AS g2s ON s2g.group_id = g2s.group_id").
		Join("LEFT JOIN supers AS relatives ON g2s.super_id = relatives.id AND g2s.super_id != s.id").
		Apply(filter).
		Group("s.id").
		Select()
	if err != nil {
		panic(err)
	}

	for i := range supersResult {

		// create the Group Names List as []string
		for _, group := range supersResult[i].Groups {
			supersResult[i].GroupsList = append(supersResult[i].GroupsList, group.Name)
		}

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
