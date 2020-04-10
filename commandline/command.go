package commandline

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-pg/pg/v9"
	db "github.com/tcarreira/superhero/models"
	"github.com/tcarreira/superhero/server"
)

// CommandLiner is an interface for methods used by ParseCommandLine
type CommandLiner interface {
	printUsage(logger *log.Logger)
	printServeUsage(logger *log.Logger)
	printAdminUsage(logger *log.Logger)
	exit(ret int)
	getArg(idx int) string
	lenArgs() int
}

// CommandLine is an instance of CommandLiner
type CommandLine struct{}

// printUsage prints usage
func (c CommandLine) printUsage(logger *log.Logger) {
	logger.Println("Usage:", filepath.Base(os.Args[0]), "COMMAND")
	logger.Println("")
	logger.Println("COMMAND:")
	logger.Println("	admin: call admin actions")
	logger.Println("	serve: start HTTP server")
}

// printServeUsage prints usage for admin sub-command
func (c CommandLine) printServeUsage(logger *log.Logger) {
	logger.Println("Usage:", filepath.Base(os.Args[0]), "serve", "COMMAND")
	logger.Println("")
	logger.Println("COMMAND:")
	logger.Println("	swagger: run server with /swagger endpoint active")
}

// printAdminUsage prints usage for admin sub-command
func (c CommandLine) printAdminUsage(logger *log.Logger) {
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

func parseCommandLine(comm CommandLiner, d *pg.DB) {
	logger := log.New(os.Stdout, "", 0)

	if comm.lenArgs() < 2 {
		comm.printUsage(logger)
		comm.exit(1)
	} else {
		switch comm.getArg(1) {

		case "serve":
			if comm.lenArgs() < 3 {
				server.RunHTTPServer(d)
			} else {
				switch comm.getArg(2) {
				case "swagger":
					server.RunHTTPServerWithSwagger(d)
				default:
					comm.printServeUsage(logger)
					comm.exit(1)
				}
			}

		case "admin":
			if comm.lenArgs() < 3 {
				comm.printAdminUsage(logger)
				comm.exit(1)
			} else {

				switch comm.getArg(2) {
				case "schema":
					db.CreateSchema(d)
				case "drop":
					db.DropSchema(d)
				case "migrate":
					db.Migrate(d)
				default:
					comm.printAdminUsage(logger)
					comm.exit(1)
				}
			}

		default:
			comm.printUsage(logger)
			comm.exit(1)
		}
	}
}

// Parse will parse command line arguments/flags and return a callable function
func Parse(db *pg.DB) {
	var comm CommandLine
	comm = CommandLine{}

	parseCommandLine(comm, db)
}
