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

	const q = `
		SELECT
			id, name, description, date_created, date_updated
		FROM station_types`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "selecting station types")
	}

	return products, nil
}

// Retrieve gets a specific StationType from the database.
func Retrieve(db *sqlx.DB, id string) (*StationTypes, error) {

	var st StationTypes

	const q = `
		SELECT
			id, name, description, date_created, date_updated
		FROM station_types
		WHERE id = $1`

	if err := db.Get(&st, q, id); err != nil {
		return nil, errors.Wrap(err, "selecting single station type")
	}

	return &st, nil
}
