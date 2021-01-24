package schema

import (
	// Third-party packages
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Including the queries directly in this file has the same pros/cons mentioned
// in seeds.go

var migrations = []darwin.Migration{
	{
		Version:	 1,
		Description: "Add station types",
		Script: `
CREATE TABLE station_type (
	id           UUID PRIMARY KEY,
	name         TEXT,
	description  TEXT,
	date_created TIMESTAMP,
	date_updated TIMESTAMP
);`,
	},
	{
		Version:     2,
		Description: "Add station",
		Script: `
CREATE TABLE station (
	id              UUID PRIMARY KEY,
	station_type_id UUID,
	name            TEXT,
	description     TEXT,
	location_x      INT,
	location_y      INT,
	date_created    TIMESTAMP,
	date_updated    TIMESTAMP,

	CONSTRAINT fk_station_type_id
		FOREIGN KEY (station_type_id)
		REFERENCES station_type(id)
		ON DELETE CASCADE
);`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {

	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
