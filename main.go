package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

// Server stores both Database and HTTP Server connectors
type Server struct {
	DB     *pg.DB
	Router *gin.Engine
}

func main() {
	s := Server{}

	s.setupRouter()

	s.setupDatabase()
	defer s.DB.Close()

	s.Router.Run()
}
