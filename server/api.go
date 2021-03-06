package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/tcarreira/superhero/docs" // swagger import side effects

	"github.com/tcarreira/superhero/models"
)

// @title Superhero API
// @version 1.0
// @description SuperHero API - Go (inspired by superheroapi.com) \n This is being made in the context of https://github.com/levpay/careers#desafio

// @contact.name Tiago Carreira
// @contact.url https://github.com/tcarreira/superhero

// @license.name MIT License
// @license.url https://raw.githubusercontent.com/tcarreira/superhero/master/LICENSE

// @BasePath /api/v1

type errorResponseJSON struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

//     _____
//    / ____|
//   | (___  _   _ _ __   ___ _ __ ___
//    \___ \| | | | '_ \ / _ \ '__/ __|
//    ____) | |_| | |_) |  __/ |  \__ \
//   |_____/ \__,_| .__/ \___|_|  |___/
//                | |
//                |_|

// SuperHandler interface for REST API for Super
type SuperHandler interface {
	SuperHeroPOSTHandler(c *gin.Context)
	SuperVilanPOSTHandler(c *gin.Context)
	SupersPOSTHandler(c *gin.Context)
	SupersGETFiltersHandler(c *gin.Context)
	SupersGETByIDHandler(c *gin.Context)
	SupersPUTHandler(c *gin.Context)
	SupersDeleteHandler(c *gin.Context)
}

// SuperAPI implements SuperHandler interface
type SuperAPI struct {
	DB     *pg.DB
	Router *gin.Engine
}

func (api *SuperAPI) handleSuperCreate(c *gin.Context, super *models.Super) {
	if _, err := super.Create(api.DB); err != nil {
		if _, ok := err.(*models.ErrorSuperAlreadyExists); ok {
			c.JSON(http.StatusConflict, errorResponseJSON{
				"Super already exists - update it instead",
				err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponseJSON{
				"Unexpected error",
				err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusCreated, super)
	}
}

func (api *SuperAPI) handleSuperBindingJSON(c *gin.Context) (*models.Super, bool) {
	super := models.Super{}

	if err := c.ShouldBindJSON(&super); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseJSON{
			"Error processing the payload",
			err.Error(),
		})
		return &super, false
	}
	return &super, true
}

type exampleSuperHeroVilanJSON struct {
	Name string `json:"name" example:"name1"`
}

// SuperHeroPOSTHandler Create new SuperHero
// ---
// @Summary Create new Super Hero
// @Description Create new Super Hero by name
// @Accept  json
// @Produce  json
// @Param super body exampleSuperHeroVilanJSON true "super hero name"
// @Success 201 {object} models.Super "Super was created"
// @Failure 409 {object} errorResponseJSON "Super already exists"
// @Failure 500 {object} errorResponseJSON "Unexpected error"
// @Router /super-hero [post]
func (api *SuperAPI) SuperHeroPOSTHandler(c *gin.Context) {

	super, ok := api.handleSuperBindingJSON(c)
	if !ok {
		return
	}
	super.Type = "HERO"

	api.handleSuperCreate(c, super)
}

// SuperVilanPOSTHandler Create new Super Vilan
// ---
// @Summary Create new Super Vilan
// @Description Create new Super Vilan by name
// @Accept  json
// @Produce  json
// @Param super body exampleSuperHeroVilanJSON true "super vilan name"
// @Success 201 {object} models.Super "Super was created"
// @Failure 409 {object} errorResponseJSON "Super already exists"
// @Router /super-vilan [post]
func (api *SuperAPI) SuperVilanPOSTHandler(c *gin.Context) {

	super, ok := api.handleSuperBindingJSON(c)
	if !ok {
		return
	}
	super.Type = "VILAN"

	api.handleSuperCreate(c, super)
}

type exampleSuperJSON struct {
	Type string `json:"type" example:"HERO"`
	Name string `json:"name" example:"name1"`
}

// SupersPOSTHandler Create new Super
// ---
// @Summary Create new Super (hero/vilan)
// @Description Create new Super
// @Accept  json
// @Produce  json
// @Param super body exampleSuperJSON true "super hero (mandatory: name and type)"
// @Success 201 {object} models.Super "Super was created"
// @Failure 409 {object} errorResponseJSON "Super already exists"
// @Router /supers [post]
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
// @Success 200 {array} models.Super "List of Supers"
// @Failure 400 {object} errorResponseJSON "Error parsing payload"
// @Router /supers [get]
func (api *SuperAPI) SupersGETFiltersHandler(c *gin.Context) {
	sFilter := models.Super{}

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
// @Success 200 {object} models.Super "Super"
// @Failure 404 {object} errorResponseJSON "Super Not Found"
// @Failure 500 {object} errorResponseJSON "Unexpected Error"
// @Router /supers/{id} [get]
func (api *SuperAPI) SupersGETByIDHandler(c *gin.Context) {
	var super *models.Super
	var err error

	super, err = super.GetByNameOrUUID(api.DB, c.Param("id"))
	if err != nil {
		if _, ok := err.(*models.ErrorSuperNotFound); ok {
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

// SupersPUTHandler Update Super
func (api *SuperAPI) SupersPUTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"Updating Super is out of scope for now",
	})
}

// SupersDeleteHandler Delete a Super @ /supers/:name...
// ---
// @Summary Delete a Super
// @Description Delete a by name or uuid
// @Produce json
// @Param id path string true "Super's Name or UUID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} errorResponseJSON "Super Not Found"
// @Failure 500 {object} errorResponseJSON "Unexpected Error"
// @Router /supers/{id} [delete]
func (api *SuperAPI) SupersDeleteHandler(c *gin.Context) {
	super := new(models.Super)
	err := super.DeleteByNameOrUUID(api.DB, c.Param("id"))

	if err != nil {
		if _, ok := err.(*models.ErrorSuperNotFound); ok {
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

//////////////////////////////////////////////////////////////
//     _____
//    / ____|
//   | |  __ _ __ ___  _   _ _ __  ___
//   | | |_ | '__/ _ \| | | | '_ \/ __|
//   | |__| | | | (_) | |_| | |_) \__ \
//    \_____|_|  \___/ \__,_| .__/|___/
//                          | |
//                          |_|
//////////////////////////////////////////////////////////////

// GroupHandler interface for REST API for Groups
type GroupHandler interface {
	GroupsPOSTHandler(c *gin.Context)
	GroupsGETHandler(c *gin.Context)
	GroupsPUTHandler(c *gin.Context)
	GroupsDeleteHandler(c *gin.Context)
}

// GroupAPI implements GroupHandler interface
type GroupAPI struct {
	DB     *pg.DB
	Router *gin.Engine
}

// GroupsPOSTHandler Create Group
// ---
// @Summary Create new Group of Supers
// @Description Create new Group of Supers
// @Accept  json
// @Produce  json
// @Param super body models.Group true "Group definition. Supers is a list os their names"
// @Success 201 {object} models.Group "Group was created"
// @Failure 409 {object} errorResponseJSON "Group name already exists"
// @Failure 500 {object} errorResponseJSON "Unexpected Error"
// @Router /groups [post]
func (api *GroupAPI) GroupsPOSTHandler(c *gin.Context) {
	group := models.Group{}

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseJSON{
			"Error processing the payload",
			err.Error(),
		})
		return
	}

	if group, err := group.Create(api.DB); err != nil {
		if _, ok := err.(*models.ErrorGroupAlreadyExists); ok {
			c.JSON(http.StatusConflict, errorResponseJSON{
				"Group already exists - update it instead",
				err.Error(),
			})
		} else if _, ok := err.(*models.ErrorGroupSuperRelation); ok {
			log.Println("Found some non-fatal errors. Will log and ignore:", err.Error())
			c.JSON(http.StatusCreated, group)
		} else {
			c.JSON(http.StatusInternalServerError, errorResponseJSON{
				"Unexpected Error",
				err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusCreated, group)
	}

}

// GroupsGETHandler Get a Group
// ---
// @Summary Get Group
// @Description Get Group by name
// @Produce json
// @Param name path string true "Group Name"
// @Success 200 {object} models.Group "Group"
// @Failure 404 {object} errorResponseJSON "Group Not Found"
// @Failure 500 {object} errorResponseJSON "Unexpected Error"
// @Router /groups/{name} [get]
func (api *GroupAPI) GroupsGETHandler(c *gin.Context) {
	var group *models.Group
	var err error

	group, err = group.GetByName(api.DB, c.Param("name"))
	if err != nil {
		if _, ok := err.(*models.ErrorGroupNotFound); ok {
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

// GroupsPUTHandler Update a Group
func (api *GroupAPI) GroupsPUTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not implemented", "out of the scope for this exercise",
	})
}

// GroupsDeleteHandler Delete a Group
func (api *GroupAPI) GroupsDeleteHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not implemented", "out of the scope for this exercise",
	})
}

//    _____             _
//   |  __ \           | |
//   | |__) |___  _   _| |_ ___  ___
//   |  _  // _ \| | | | __/ _ \/ __|
//   | | \ \ (_) | |_| | ||  __/\__ \
//   |_|  \_\___/ \__,_|\__\___||___/
//
//

func setRoutes(r *gin.Engine, db *pg.DB) *gin.Engine {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	v1 := r.Group("/api/v1")
	{
		// Supers
		{
			api := SuperAPI{
				DB:     db,
				Router: r,
			}

			v1.POST("/super-hero", api.SuperHeroPOSTHandler)
			v1.POST("/super-vilan", api.SuperVilanPOSTHandler)

			supers := v1.Group("/supers")
			{
				supers.POST("", api.SupersPOSTHandler)
				supers.GET("", api.SupersGETFiltersHandler)
				supers.GET("/:id", api.SupersGETByIDHandler)
				supers.PUT("/:id", api.SupersPUTHandler)
				supers.DELETE("/:id", api.SupersDeleteHandler)
			}
		}

		groups := v1.Group("/groups")
		{
			api := GroupAPI{
				DB:     db,
				Router: r,
			}

			groups.POST("/", api.GroupsPOSTHandler)
			groups.GET("/:name", api.GroupsGETHandler)
			groups.PUT("/:name", api.GroupsPUTHandler)
			groups.DELETE("/:name", api.GroupsDeleteHandler)
		}
	}

	return r
}

//     _____
//    / ____|
//   | (___   ___ _ ____   _____ _ __
//    \___ \ / _ \ '__\ \ / / _ \ '__|
//    ____) |  __/ |   \ V /  __/ |
//   |_____/ \___|_|    \_/ \___|_|
//
//

// SetupRouter setup a default gin.Engine and setup Routes but do not run
func SetupRouter(db *pg.DB) *gin.Engine {
	r := gin.Default()
	r = setRoutes(r, db)

	return r
}

// RunHTTPServer setup routes and start http server
func RunHTTPServer(db *pg.DB) {
	r := SetupRouter(db)
	r.Run()
}

// RunHTTPServerWithSwagger setup routes (with /swagger) and start http server
func RunHTTPServerWithSwagger(db *pg.DB) {
	r := SetupRouter(db)
	// s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger/*any", ginSwagger.CustomWrapHandler(
		&ginSwagger.Config{
			URL: "doc.json",
		},
		swaggerFiles.Handler,
	))

	r.Run()
}
