package database

import (
	"net/url"

	// Third-party packages
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // https://pkg.go.dev/github.com/lib/pq
)

// Config is the required properties to use the database.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

/**
 * Open database connection.
 *
 * https://golang.org/pkg/database/sql/
 * Supported drivers: https://github.com/golang/go/wiki/SQLDrivers
 */
func Open(cfg Config) (*sqlx.DB, error) {

	// Define SSL mode.
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
