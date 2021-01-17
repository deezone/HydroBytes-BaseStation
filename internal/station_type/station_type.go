package station_types

import (
	// Core packages
	"context"       // https://golang.org/pkg/context
	"database/sql"
	"time"

	// Third-party packages
	"github.com/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific StationTypes is requested but does not exist.
	ErrNotFound = errors.New("station type not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper UUID format")
)

// Create adds a StationType to the database. It returns the created StationTypes with
// fields like ID and DateCreated populated.
func Create(ctx context.Context, db *sqlx.DB, nst NewStationTypes, now time.Time) (*StationTypes, error) {
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

	_, err := db.ExecContext(ctx, q,
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

// List gets all StationTypes from the database.
func List(ctx context.Context, db *sqlx.DB) ([]StationTypes, error) {
	station_types := []StationTypes{}

	const q = `
		SELECT
			id, name, description, date_created, date_updated
		FROM station_types`

	if err := db.SelectContext(ctx, &station_types, q); err != nil {
		return nil, errors.Wrap(err, "selecting station types")
	}

	return station_types, nil
}

// Retrieve gets a specific StationType from the database.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*StationTypes, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var st StationTypes

	const q = `
		SELECT
			id, name, description, date_created, date_updated
		FROM station_types
		WHERE id = $1`

	if err := db.GetContext(ctx, &st, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single station type")
	}

	return &st, nil
}
