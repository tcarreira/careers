package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setRoutes(r *gin.Engine) *gin.Engine {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	v1 := r.Group("/api/v1")
	{
		api := new(SuperAPI)

		v1.POST("/supers", api.supersPOSTHandler)
		v1.GET("/supers", api.supersGETHandler)
		v1.GET("/supers/:id", api.supersGETHandler)
		v1.PUT("/supers/:id", api.supersPUTHandler)
		v1.DELETE("/supers/:id", api.supersDeleteHandler)
	}

	return r
}

func (s *Server) setupRouter() *Server {
	s.Router = gin.Default()

	s.Router = setRoutes(s.Router)

	return s
}

func (s *Server) runHTTPServer() {
	s.setupRouter()
	s.Router.Run()
}
