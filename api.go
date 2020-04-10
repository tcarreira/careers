package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/tcarreira/superhero/docs"
)

// @title Superhero API
// @version 1.0
// @description SuperHero API - Go (inspired by superheroapi.com) \n This is being made in the context of https://github.com/levpay/careers#desafio

// @contact.name Tiago Carreira
// @contact.url https://github.com/tcarreira/superhero

// @license.name MIT License
// @license.url https://raw.githubusercontent.com/tcarreira/superhero/master/LICENSE

// @BasePath /api/v1

// SuperHandler interface for REST API for Super
type SuperHandler interface {
	supersPOSTHandler(c *gin.Context)
	supersGETHandler(c *gin.Context)
	supersPUTHandler(c *gin.Context)
	supersDeleteHandler(c *gin.Context)
}

// SuperAPI implements SuperHandler interface
type SuperAPI struct {
	DB     *pg.DB
	Router *gin.Engine
}

type errorResponseJSON struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (api *SuperAPI) handleSuperCreate(c *gin.Context, super *Super) {
	if _, err := super.Create(api.DB); err != nil {
		if _, ok := err.(*errorSuperAlreadyExists); ok {
			c.JSON(http.StatusConflict, errorResponseJSON{
				"Super already exists - update it instead",
				err.Error(),
			})
		} else {
			panic(err)
		}
	} else {
		c.JSON(http.StatusCreated, super)
	}
}

func (api *SuperAPI) handleSuperBindingJSON(c *gin.Context) (*Super, bool) {
	super := Super{}

	if err := c.ShouldBindJSON(&super); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseJSON{
			"Error processing the payload",
			err.Error(),
		})
		return &super, false
	}
	return &super, true
}

func (api *SuperAPI) superHeroPOSTHandler(c *gin.Context) {

	super, ok := api.handleSuperBindingJSON(c)
	if !ok {
		return
	}
	super.Type = "HERO"

	api.handleSuperCreate(c, super)
}

func (api *SuperAPI) superVilanPOSTHandler(c *gin.Context) {

	super, ok := api.handleSuperBindingJSON(c)
	if !ok {
		return
	}
	super.Type = "VILAN"

	api.handleSuperCreate(c, super)
}

type exampleSuperHeroVilanJSON struct {
	Type string `json:"type" example:"HERO"`
	Name string `json:"name" example:"name1"`
}

// SupersPOSTHandler Create new Super
// ---
// @Summary Create new Super (hero/vilan)
// @Description Create new Super
// @Accept  json
// @Produce  json
// @Param super body exampleSuperHeroVilanJSON true "super hero name"
// @Success 201 {object} Super "Super was created"
// @Failure 409 {object} errorResponseJSON "Super already exists"
// @Router /super-hero [post]
func (api *SuperAPI) SupersPOSTHandler(c *gin.Context) {

	super, ok := api.handleSuperBindingJSON(c)
	if !ok {
		return
	}

	api.handleSuperCreate(c, super)
}

// SupersGETFiltersHandler get list of Super @ /supers?type=hero...
// ---
// @Summary Get list of Supers
// @Description Get list of Supers by filtering by name, uuid or type
// @Produce json
// @Param name query string false "Super(hero/vilan) Name (case-sensitive)"
// @Param uuid query string false "Super(hero/vilan) UUID (case-insensitive)"
// @Param type query string false "Super(hero/vilan) Type (HERO / VILAN) (case-insensitive)"
// @Success 200 {array} Super "List of Supers"
// @Failure 400 {object} errorResponseJSON "Error parsing payload"
// @Router /supers [get]
func (api *SuperAPI) SupersGETFiltersHandler(c *gin.Context) {
	sFilter := Super{}

	if err := c.ShouldBind(&sFilter); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseJSON{
			"Could not process Payload (query parameters)",
			err.Error(),
		})
	} else {
		results := sFilter.ReadAll(api.DB)
		c.JSON(http.StatusOK, results)
	}
}

// SupersGETByIDHandler a Super @ /supers/:id...
// ---
// @Summary Get Super
// @Description Get a Super by name or uuid
// @Produce json
// @Param id path string true "Super's Name or UUID"
// @Success 200 {object} Super "Super"
// @Failure 404 {object} errorResponseJSON "Super Not Found"
// @Failure 500 {object} errorResponseJSON "Unexpected Error"
// @Router /supers/{id} [get]
func (api *SuperAPI) SupersGETByIDHandler(c *gin.Context) {
	var super *Super
	var err error

	super, err = super.getByNameOrUUID(api.DB, c.Param("id"))
	if err != nil {
		if _, ok := err.(*errorSuperNotFound); ok {
			c.JSON(http.StatusNotFound, errorResponseJSON{
				"No Super was found",
				err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponseJSON{
				"Internal Server Error",
				err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, super)
}

func (api *SuperAPI) supersPUTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"Updating Super is out of scope for now",
	})
}

func (api *SuperAPI) supersDeleteHandler(c *gin.Context) {
	super := new(Super)
	err := super.DeleteByNameOrUUID(api.DB, c.Param("id"))

	if err != nil {
		if _, ok := err.(*errorSuperNotFound); ok {
			// nothing was deleted - return 404
			c.JSON(http.StatusNotFound, errorResponseJSON{
				"Super Not Found",
				err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponseJSON{
				"Unexpected Error",
				err.Error(),
			})
		}
	} else {
		// Deleted. No Content is needed
		c.Status(http.StatusNoContent)
	}
}

func (api *SuperAPI) groupsPOSTHandler(c *gin.Context) {
	group := Group{}

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseJSON{
			"Error processing the payload",
			err.Error(),
		})
		return
	}

	if group, err := group.Create(api.DB); err != nil {
		if _, ok := err.(*errorGroupAlreadyExists); ok {
			c.JSON(http.StatusConflict, errorResponseJSON{
				"Group already exists - update it instead",
				err.Error(),
			})
		} else if _, ok := err.(*errorGroupSuperRelation); ok {
			log.Println("Found some non-fatal errors. Will log and ignore:", err.Error())
		} else {
			panic(err)
		}
	} else {
		c.JSON(http.StatusCreated, group)
	}

}

func (api *SuperAPI) groupsGETHandler(c *gin.Context) {
	var group *Group
	var err error

	group, err = group.GetByName(api.DB, c.Param("name"))
	if err != nil {
		if _, ok := err.(*errorGroupNotFound); ok {
			c.JSON(http.StatusNotFound, errorResponseJSON{
				"Group not found", err.Error(),
			})
		} else {
			panic(err)
		}
	} else {
		c.JSON(http.StatusOK, group)
	}
}

func (api *SuperAPI) groupsPUTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not implemented", "out of the scope for this exercise",
	})
}

func (api *SuperAPI) groupsDeleteHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not implemented", "out of the scope for this exercise",
	})
}

//  Routes

func setRoutes(s *Server) *gin.Engine {

	s.Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	v1 := s.Router.Group("/api/v1")
	{
		api := SuperAPI{
			DB:     s.DB,
			Router: s.Router,
		}

		v1.POST("/super-hero", api.superHeroPOSTHandler)
		v1.POST("/super-vilan", api.superVilanPOSTHandler)

		v1.POST("/supers", api.SupersPOSTHandler)
		v1.GET("/supers", api.SupersGETFiltersHandler)
		v1.GET("/supers/:id", api.SupersGETByIDHandler)
		v1.PUT("/supers/:id", api.supersPUTHandler)
		v1.DELETE("/supers/:id", api.supersDeleteHandler)

		v1.POST("/groups", api.groupsPOSTHandler)
		v1.GET("/groups/:name", api.groupsGETHandler)
		v1.PUT("/groups/:name", api.groupsPUTHandler)
		v1.DELETE("/groups/:name", api.groupsDeleteHandler)
	}

	return s.Router
}

// run server

func (s *Server) setupRouter() *Server {
	s.Router = gin.Default()

	s.Router = setRoutes(s)

	return s
}

func (s *Server) runHTTPServer() {
	s.setupRouter()
	s.Router.Run()
}

func (s *Server) runHTTPServerWithSwagger() {
	s.setupRouter()
	// s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.Router.GET("/swagger/*any", ginSwagger.CustomWrapHandler(
		&ginSwagger.Config{
			URL: "doc.json",
		},
		swaggerFiles.Handler,
	))

	s.Router.Run()
}
