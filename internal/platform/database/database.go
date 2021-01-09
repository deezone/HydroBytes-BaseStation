package database

import (
	"net/url"

	// Third-party packages
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // https://pkg.go.dev/github.com/lib/pq
)

/**
 * Open database connection
 *
 * https://golang.org/pkg/database/sql/
 * Supported drivers: https://github.com/golang/go/wiki/SQLDrivers
 */
func Open() (*sqlx.DB, error) {
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
