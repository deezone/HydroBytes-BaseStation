package station_type

import (
	// Core packages
	"context"
	"database/sql"
	"time"

	// Third-party packages
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AddStation adds a station of a specific StationType.
func AddStation(ctx context.Context, db *sqlx.DB, ns NewStation, stationTypeID string, now time.Time) (*Station, error) {
	s := Station{
		Id:            uuid.New().String(),
		StationTypeId: stationTypeID,
		Name:          ns.Name,
		Description:   ns.Description,
		LocationX:     ns.LocationX,
		LocationY:     ns.LocationY,
		DateCreated:   now,
		DateUpdated:   now,
	}

	const q = `INSERT INTO station
		(id, station_type_id, name, description, location_x, location_y, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.ExecContext(ctx, q,
		s.Id,
		s.StationTypeId,
		s.Name,
		s.Description,
		s.LocationX,
		s.LocationY,
		s.DateCreated,
		s.DateUpdated,
	)
	if err != nil {
		return nil, errors.Wrap(err, "inserting station")
	}

	return &s, nil
}

// AdjustStation modifies data about a Station. It will error if the specified ID is
// invalid or does not reference an existing Station.
func AdjustStation(ctx context.Context, db *sqlx.DB, id string, update UpdateStation, now time.Time) error {
	s, err := RetrieveStation(ctx, db, id)
	if err != nil {
		return err
	}

	if update.Name != nil {
		s.Name = *update.Name
	}
	if update.Description != nil {
		s.Description = *update.Description
	}
	if update.LocationX != nil {
		s.LocationX = *update.LocationX
	}
	if update.LocationY != nil {
		s.LocationY = *update.LocationY
	}
	s.DateUpdated = now

	const q = `UPDATE station SET
		"name" = $2,
		"description" = $3,
        "location_x" = $4,
        "location_x" = $5,
		WHERE id = $1`
	_, err = db.ExecContext(ctx, q, id,
		s.Name,
		s.Description,
		s.LocationX,
		s.LocationY,
	)
	if err != nil {
		return errors.Wrap(err, "updating station")
	}

	return nil
}

// ListStations gives all Stations for a StationType.
func ListStations(ctx context.Context, db *sqlx.DB, stationTypeID string) ([]Station, error) {
	stations := []Station{}

	const q = `
      SELECT
        id,
        station_type_id,
        name, description,
        location_x,
        location_y,
        date_created,
        date_updated
      FROM station
      WHERE station_type_id = $1`

	if err := db.SelectContext(ctx, &stations, q, stationTypeID); err != nil {
		return nil, errors.Wrap(err, "selecting stations")
	}

	return stations, nil
}

// Retrieve gets a specific Station from the database.
func RetrieveStation(ctx context.Context, db *sqlx.DB, id string) (*Station, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var s Station

	const q = `
		SELECT
			id,
		    station_type_id,
			name,
			description,
			location_x,
			location_y,
			date_created,
			date_updated
		FROM station
		WHERE id = $1`

	if err := db.GetContext(ctx, &s, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single station")
	}

	return &s, nil
}
