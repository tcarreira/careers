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
				err.Error(),
				"",
			})
		} else {
			panic(err)
		}
	} else {
		c.JSON(http.StatusCreated, super)
	}

}

func (api *SuperAPI) supersGETHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"supersGETHandler WIP",
	})
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
