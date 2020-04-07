package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuperHandler interface for REST API for Super
type SuperHandler interface {
	supersPOSTHandler(c *gin.Context)
	supersGETHandler(c *gin.Context)
	supersPUTHandler(c *gin.Context)
	supersDeleteHandler(c *gin.Context)
}

// SuperAPI implements SuperHandler interface
type SuperAPI struct{}

type errorResponseJSON struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (api *SuperAPI) supersPOSTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"supersPOSTHandler WIP",
	})
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
