package station_types

import (
	// Core packages
	"github.com/pkg/errors"

	// Third party packages
	"github.com/jmoiron/sqlx"
)

// List gets all Products from the database.
func List(db *sqlx.DB) ([]StationTypes, error) {
	products := []StationTypes{}

	const q = `SELECT id, name, description, date_created, date_updated FROM station_types`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "selecting station types")
	}

	return products, nil
}
