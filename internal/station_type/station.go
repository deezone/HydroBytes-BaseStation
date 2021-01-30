package station_type

import (
	// Core packages
	"context"
	"database/sql"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/auth"

	// Third-party packages
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrStationNotFound is used when a specific Station is requested but does not exist.
	ErrStationNotFound = errors.New("station not found")

	// ErrForbidden occurs when an account tries to do something that is forbidden to
	// it according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// AddStation adds a station of a specific StationType.
func AddStation(ctx context.Context, db *sqlx.DB, account auth.Claims, ns NewStation, stationTypeID string, now time.Time) (*Station, error) {
	s := Station{
		Id:            uuid.New().String(),
		StationTypeId: stationTypeID,
		AccountId:     account.Subject,
		Name:          ns.Name,
		Description:   ns.Description,
		LocationX:     ns.LocationX,
		LocationY:     ns.LocationY,
		DateCreated:   now.UTC(),
		DateUpdated:   now.UTC(),
	}

	const q = `INSERT INTO station
		(id, station_type_id, account_id, name, description, location_x, location_y, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := db.ExecContext(ctx, q,
		s.Id,
		s.StationTypeId,
		s.AccountId,
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
func AdjustStation(ctx context.Context, db *sqlx.DB, account auth.Claims, id string, update UpdateStation, now time.Time) error {
	s, err := RetrieveStation(ctx, db, id)
	if err != nil {
		return err
	}

	// If the account attempt to access does not have the admin role ...
	// or the owner of the station ...
	// then restrict access
	if !account.HasRole(auth.RoleAdmin) && s.AccountId != account.Subject {
		return ErrForbidden
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
        "location_y" = $5,
        "date_updated" = $6
		WHERE id = $1`
	_, err = db.ExecContext(ctx, q, id,
		s.Name,
		s.Description,
		s.LocationX,
		s.LocationY,
		s.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating station")
	}

	return nil
}

// DeleteStation removes the station identified by a given ID.
func DeleteStation(ctx context.Context, db *sqlx.DB, id string) error {
	// Validate id is a valid uuid
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM station WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting station %s", id)
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
        account_id,
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
            account_id,
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
            return nil, ErrStationNotFound
        }

        return nil, errors.Wrap(err, "selecting single station")
    }

    return &s, nil
}
