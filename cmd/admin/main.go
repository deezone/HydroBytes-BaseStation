package main

import (
	"flag"
	"log" // https://golang.org/pkg/log/
	"net/url"
	"os"

	// Applcation packages
	"github.com/deezone/HydroBytes-BaseStation/internal/schema"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/database"

	// Third-party packages
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // https://pkg.go.dev/github.com/lib/pq
)

// Main entry point for command line functionality.
func main() {

	db, err := database.Open()
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	flag.Parse()

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

/**
 * https://golang.org/pkg/database/sql/
 * Supported drivers: https://github.com/golang/go/wiki/SQLDrivers
 */
func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
