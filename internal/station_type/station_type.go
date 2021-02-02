package station_type

import (
	// Core packages
	"context"       // https://golang.org/pkg/context
	"database/sql"
	"time"

	// Third-party packages
	"github.com/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific StationType is requested but does not exist.
	ErrNotFound = errors.New("station type not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper UUID format")
)

// Create adds a StationType to the database. It returns the created StationType with
// fields like ID and DateCreated populated.
func Create(ctx context.Context, db *sqlx.DB, nst NewStationType, now time.Time) (*StationType, error) {

	ctx, span := trace.StartSpan(ctx, "station_type.Create")
	defer span.End()

	st := StationType{
		Id:          uuid.New().String(),
		Name:        nst.Name,
		Description: nst.Description,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
		INSERT INTO station_type
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

// Delete removes the station type identified by a given ID.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {

	ctx, span := trace.StartSpan(ctx, "station_type.Delete")
	defer span.End()

	// Validate id is a valid uuid
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM station_type WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting station type %s", id)
	}

	return nil
}

// List gets all StationType from the database.
func List(ctx context.Context, db *sqlx.DB) ([]StationType, error) {

	ctx, span := trace.StartSpan(ctx, "station_type.List")
	defer span.End()

	station_type := []StationType{}

	const q = `
		SELECT
			station_type.id,
			station_type.name,
			station_type.description,
			COUNT(station.id) AS stations,
			station_type.date_created,
			station_type.date_updated
		FROM station_type
		  LEFT JOIN station ON station_type.id = station.station_type_id
		GROUP BY station_type.id`

	if err := db.SelectContext(ctx, &station_type, q); err != nil {
		return nil, errors.Wrap(err, "selecting station types")
	}

	return station_type, nil
}

// Get finds the product identified by a given ID.
func Get(ctx context.Context, db *sqlx.DB, id string) (*StationType, error) {

	ctx, span := trace.StartSpan(ctx, "station_type.Get")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var st StationType

	const q = `
		SELECT
			station_type.id,
			station_type.name,
			station_type.description,
			COUNT(station.id) AS stations,
			station_type.date_created,
			station_type.date_updated
		FROM station_type
		  LEFT JOIN station ON station_type.id = station.station_type_id
		WHERE station_type.id = $1
		GROUP BY station_type.id`

	if err := db.GetContext(ctx, &st, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single station type")
	}

	return &st, nil
}

// Update modifies data about a StationType. It will error if the specified ID is
// invalid or does not reference an existing StationType.
func Update(ctx context.Context, db *sqlx.DB, id string, update UpdateStationType, now time.Time) error {

	ctx, span := trace.StartSpan(ctx, "station_type.Update")
	defer span.End()

	st, err := Get(ctx, db, id)
	if err != nil {
		return err
	}

	if update.Name != nil {
		st.Name = *update.Name
	}
	if update.Description != nil {
		st.Description = *update.Description
	}
	st.DateUpdated = now

	const q = `UPDATE station_type SET
		"name" = $2,
		"description" = $3,
		"date_updated" = $4
		WHERE id = $1`
	_, err = db.ExecContext(ctx, q, id,
		st.Name,
		st.Description,
		st.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating station tyoe")
	}

	return nil
}
