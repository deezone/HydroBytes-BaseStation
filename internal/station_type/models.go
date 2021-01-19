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

// UpdateStationType defines what information may be provided to modify an
// existing StationType. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that were not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateStationType struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
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
	LocationX     int       `db:"location_x"      json:"location_x" validate:"required,gte=0"`
	LocationY     int       `db:"location_y"      json:"location_y" validate:"required,gte=0"`
}

// UpdateStation defines what information may be provided to modify an
// existing Station. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that were not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateStation struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	LocationX    *int    `json:"location_x" validate:"omitempty,gte=0"`
	LocationY    *int    `json:"location_y" validate:"omitempty,gte=0"`
}
