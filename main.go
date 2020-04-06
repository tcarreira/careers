package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setRoutes(r *gin.Engine) *gin.Engine {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	return r
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	return setRoutes(r)
}

func main() {
	setupRouter().Run()
}
