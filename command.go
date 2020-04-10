package main

import (
	"log"
	"os"
	"path/filepath"
)

// CommandLiner is an interface for methods used by ParseCommandLine
type CommandLiner interface {
	usagePrint(logger *log.Logger)
	serveUsagePrint(logger *log.Logger)
	adminUsagePrint(logger *log.Logger)
	exit(ret int)
	getArg(idx int) string
	lenArgs() int
}

// CommandLine is an instance of CommandLiner
type CommandLine struct{}

// usagePrint prints usage
func (c CommandLine) usagePrint(logger *log.Logger) {
	logger.Println("Usage:", filepath.Base(os.Args[0]), "COMMAND")
	logger.Println("")
	logger.Println("COMMAND:")
	logger.Println("	admin: call admin actions")
	logger.Println("	serve: start HTTP server")
}

// serveUsagePrint prints usage for admin sub-command
func (c CommandLine) serveUsagePrint(logger *log.Logger) {
	logger.Println("Usage:", filepath.Base(os.Args[0]), "serve", "COMMAND")
	logger.Println("")
	logger.Println("COMMAND:")
	logger.Println("	swagger: run server with /swagger endpoint active")
}

// adminUsagePrint prints usage for admin sub-command
func (c CommandLine) adminUsagePrint(logger *log.Logger) {
	logger.Println("Usage:", filepath.Base(os.Args[0]), "admin", "COMMAND")
	logger.Println("")
	logger.Println("COMMAND:")
	logger.Println("	schema: create database schema")
	logger.Println("	migrate: perform database migrations (not implemented)")
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
	logger := log.New(os.Stdout, "", 0)

	if comm.lenArgs() < 2 {
		comm.usagePrint(logger)
		comm.exit(1)
	} else {
		switch comm.getArg(1) {

		case "serve":
			if comm.lenArgs() < 3 {
				s.runHTTPServer()
			} else {
				switch comm.getArg(2) {
				case "swagger":
					s.runHTTPServerWithSwagger()
				default:
					comm.serveUsagePrint(logger)
					comm.exit(1)
				}
			}

		case "admin":
			if comm.lenArgs() < 3 {
				comm.adminUsagePrint(logger)
				comm.exit(1)
			} else {

				switch comm.getArg(2) {
				case "schema":
					s.DBCreateSchema()
				case "migrate":
					s.DBMigrate()
				default:
					comm.adminUsagePrint(logger)
					comm.exit(1)
				}
			}

		default:
			comm.usagePrint(logger)
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
