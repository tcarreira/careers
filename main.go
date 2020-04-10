package main

import (
	"github.com/tcarreira/superhero/commandline"
	db "github.com/tcarreira/superhero/models"
)

func main() {
	db := db.SetupDatabase()
	defer db.Close()

	commandline.Parse(db)
}
