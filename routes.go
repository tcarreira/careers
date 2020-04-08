package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

		v1.POST("/supers", api.supersPOSTHandler)
		v1.GET("/supers", api.supersGETHandler)
		v1.GET("/supers/:id", api.supersGETHandler)
		v1.PUT("/supers/:id", api.supersPUTHandler)
		v1.DELETE("/supers/:id", api.supersDeleteHandler)

		v1.POST("/groups", api.groupsPOSTHandler)
		v1.GET("/groups/:name", api.groupsGETHandler)
		v1.PUT("/groups/:name", api.groupsPUTHandler)
		v1.DELETE("/groups/:name", api.groupsDeleteHandler)
	}

	return s.Router
}

func (s *Server) setupRouter() *Server {
	s.Router = gin.Default()

	s.Router = setRoutes(s)

	return s
}

func (s *Server) runHTTPServer() {
	s.setupRouter()
	s.Router.Run()
}
