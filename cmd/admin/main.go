package main

import (
	"flag"
	"log" // https://golang.org/pkg/log/
	"os"

	"github.com/deezone/HydroBytes-BaseStation/internal/platform/database"
	// Internal applcation packages
	"github.com/deezone/HydroBytes-BaseStation/internal/schema"
)

// Main entry point for command line functionality.
func main() {

	flag.Parse()

	db, err := database.Open()
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Println("error applying migrations", err)
			os.Exit(1)
		}
		log.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Println("error seeding database", err)
			os.Exit(1)
		}
		log.Println("Seed data complete")
		return
	}
}
