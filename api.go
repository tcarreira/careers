package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

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

func (api *SuperAPI) supersPOSTHandler(c *gin.Context) {
	super := Super{}

	if err := c.ShouldBindJSON(&super); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseJSON{
			"Error processing the payload",
			err.Error(),
		})
		return
	}

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

// supersGETHandler handles GET Requests.
// Valid requests:
// 	GET /supers/name
// 	GET /supers/uuid
func (api *SuperAPI) supersGETHandler(c *gin.Context) {

	if c.Param("id") != "" { // Searching for a specific user (by name or uuid)
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

	} else {
		c.JSON(http.StatusNotImplemented, errorResponseJSON{
			"Not Implemented",
			"supersGETHandler WIP",
		})
	}
}

func (api *SuperAPI) supersPUTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"Updating Super is out of scope for now",
	})
}

func (api *SuperAPI) supersDeleteHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"supersDeleteHandler WIP",
	})
}
