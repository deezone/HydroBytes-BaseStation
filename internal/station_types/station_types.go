package station_types

import (
	// Core packages
	"time"

	// Third party packages
	"github.com/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Create adds a Product to the database. It returns the created Product with
// fields like ID and DateCreated populated..
func Create(db *sqlx.DB, nst NewStationTypes, now time.Time) (*StationTypes, error) {
	st := StationTypes{
		Id:          uuid.New().String(),
		Name:        nst.Name,
		Description: nst.Description,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
		INSERT INTO station_types
		  (id, name, description, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(q,
		st.Id,
		st.Name,
		st.Description,
		st.DateCreated,
		st.DateUpdated,
	)

	if err != nil {
		return nil, errors.Wrap(err, "inserting station tyoe")
	}

	return &st, nil
}

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
