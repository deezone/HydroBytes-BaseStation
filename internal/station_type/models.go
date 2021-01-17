package station_type

import (
	// Core packages
	"time"
)

/**
 * StationType is a type of station in the automated garden system.
 *
 * Note: use of "stuct tags" (ex: `json:"id"`) to manage the names of properties to be lowercase and snake_case. Due to
 * the use of case for visibility in Go "id" rather and "Id" would result in the value being excluded in the JSON
 * response as the encoding/json package is external to this package.
 *
 * Note: the use of db:"id" allows renaming to map to the column used in the database
 */
type StationType struct {
	Id          string    `db:"id"           json:"id"`
	Name        string    `db:"name"         json:"name"`
	Description string    `db:"description"  json:"description"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewStationType is what we require from clients when adding a StationType.
type NewStationType struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
