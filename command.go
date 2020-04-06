package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// printUsage prints usage
func printUsage() {
	fmt.Println("Usage:", filepath.Base(os.Args[0]), "COMMAND")
	fmt.Println("")
	fmt.Println("COMMAND:")
	fmt.Println("	admin: call admin actions")
	fmt.Println("	serve: start HTTP server")
}

// printAdminUsage prints usage for admin sub-command
func printAdminUsage() {
	fmt.Println("Usage:", filepath.Base(os.Args[0]), "admin", "COMMAND")
	fmt.Println("")
	fmt.Println("COMMAND:")
	fmt.Println("	schema: create database schema")
	fmt.Println("	migrate: perform database migrations (not implemented)")
}

// parseCommandLine will parse command line arguments/flags and return a callable function
func (s *Server) parseCommandLine() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {

	case "serve":
		s.runHTTPServer()

	case "admin":
		if len(os.Args) < 3 {
			printAdminUsage()
			os.Exit(1)
		}

		switch os.Args[2] {
		case "schema":
			s.DBCreateSchema()
		case "migrate":
			s.DBMigrate()
		default:
			printAdminUsage()
			os.Exit(1)
		}

	default:
		printUsage()
		os.Exit(1)

	}
}
