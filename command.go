package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// CommandLiner is an interface for methods used by ParseCommandLine
type CommandLiner interface {
	usagePrint()
	adminUsagePrint()
	exit(ret int)
	getArg(idx int) string
	lenArgs() int
}

// CommandLine is an instance of CommandLiner
type CommandLine struct{}

// usagePrint prints usage
func (c CommandLine) usagePrint() {
	fmt.Println("Usage:", filepath.Base(os.Args[0]), "COMMAND")
	fmt.Println("")
	fmt.Println("COMMAND:")
	fmt.Println("	admin: call admin actions")
	fmt.Println("	serve: start HTTP server")
}

// adminUsagePrint prints usage for admin sub-command
func (c CommandLine) adminUsagePrint() {
	fmt.Println("Usage:", filepath.Base(os.Args[0]), "admin", "COMMAND")
	fmt.Println("")
	fmt.Println("COMMAND:")
	fmt.Println("	schema: create database schema")
	fmt.Println("	migrate: perform database migrations (not implemented)")
}

// exit just calls os.Exit()
func (c CommandLine) exit(ret int) {
	os.Exit(ret)
}

// getArg gets os.Args[idx]
func (c CommandLine) getArg(idx int) string {
	return os.Args[idx]
}

// getArg gets os.Args[idx]
func (c CommandLine) lenArgs() int {
	return len(os.Args)
}

func parseCommandLine(comm CommandLiner, s *Server) {

	if comm.lenArgs() < 2 {
		comm.usagePrint()
		comm.exit(1)
	} else {
		switch comm.getArg(1) {

		case "serve":
			s.runHTTPServer()

		case "admin":
			if comm.lenArgs() < 3 {
				comm.adminUsagePrint()
				comm.exit(1)
			} else {

				switch comm.getArg(2) {
				case "schema":
					s.DBCreateSchema()
				case "migrate":
					s.DBMigrate()
				default:
					comm.adminUsagePrint()
					comm.exit(1)
				}
			}

		default:
			comm.usagePrint()
			comm.exit(1)
		}
	}
}

// ParseCommandLine will parse command line arguments/flags and return a callable function
func (s *Server) ParseCommandLine() {
	var comm CommandLine
	comm = CommandLine{}

	parseCommandLine(comm, s)
}
