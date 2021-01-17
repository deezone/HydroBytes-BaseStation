package station_type

import (
	// Core packages
	"context"
	"time"

	// Third-party packages
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AddSale records a sales transaction for a single Product.
func AddStation(ctx context.Context, db *sqlx.DB, ns NewStation, stationTypeID string, now time.Time) (*Station, error) {
	s := Station{
		Id:            uuid.New().String(),
		StationTypeId: stationTypeID,
		Name:          ns.Name,
		Description:   ns.Description,
		DateCreated:   now,
		DateUpdated:   now,
	}

	const q = `INSERT INTO station
		(station_id, station_type_id, name, description, location_x, location_y, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.ExecContext(ctx, q,
		s.Id, s.StationTypeId,
		s.Name, s.Description,
		s.LocationX, s.LocationY,
		s.DateCreated, s.DateUpdated,
	)
	if err != nil {
		return nil, errors.Wrap(err, "inserting sale")
	}

	return &s, nil
}

// ListStation gives all Stations for a StationType.
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
