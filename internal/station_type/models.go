package station_type

import (
	// Core packages
	"time"
)

/**
 * StationType is a type of station in the automated garden system.
 *
 * Note: use of "struct tags" (ex: `json:"id"`) to manage the names of properties to be lowercase and snake_case. Due to
 * the use of case for visibility in Go "id" rather and "Id" would result in the value being excluded in the JSON
 * response as the encoding/json package is external to this package.
 *
 * Note: the use of db:"id" allows renaming to map to the column used in the database
 */
type StationType struct {
	Id          string    `db:"id"           json:"id"`
	Name        string    `db:"name"         json:"name"`
	Description string    `db:"description"  json:"description"`
	Stations    int       `db:"stations"     json:"stations"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewStationType is what we require from clients when adding a StationType.
type NewStationType struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
}

// Station is a station that is defined as one of the station types in StationType.
type Station struct {
	Id            string    `db:"id"              json:"id"`
	StationTypeId string    `db:"station_type_id" json:"station_type_id"`
	Name          string    `db:"name"            json:"name"`
	Description   string    `db:"description"     json:"description"`
	LocationX     int       `db:"location_x"      json:"location_x"`
	LocationY     int       `db:"location_y"      json:"location_y"`
	DateCreated   time.Time `db:"date_created"    json:"date_created"`
	DateUpdated   time.Time `db:"date_updated"    json:"date_updated"`
}

// NewStation is a what we require from clients when adding a BaseStation.
type NewStation struct {
	Name          string    `db:"name"            json:"name" validate:"required"`
	Description   string    `db:"description"     json:"description"`
	LocationX     int       `db:"location_x"      json:"location_x" validate:"gte=0"`
	LocationY     int       `db:"location_y"      json:"location_y" validate:"gte=0"`
}
