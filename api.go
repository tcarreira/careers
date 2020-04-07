package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorResponseJSON struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func supersPOSTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"supersPOSTHandler WIP",
	})
}

func supersGETHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"supersGETHandler WIP",
	})
}

func supersPUTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"Updating Super is out of scope for now",
	})
}

func supersDeleteHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"supersDeleteHandler WIP",
	})
}
