package schema

import (
	"github.com/jmoiron/sqlx"
)

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Using a constant in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.

const seeds = `
DELETE FROM station;
DELETE FROM station_type;

INSERT INTO station_type (id, name, description, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'Base', 'Coordinator for all station types - monitor, command and control. Access point to public Intenet.', '2021-01-01 00:00:01.000001+00', '2021-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'Water', 'Management of water resources. Controls water levels in resavour and impliments irrigation.', '2021-01-01 00:00:02.000001+00', '2021-01-01 00:00:02.000001+00'),
    ('5c86bbaa-4ef8-11eb-ae93-0242ac130002', 'Plant', 'Monitors and reports plant health.', '2021-01-01 00:00:03.000001+00', '2021-01-01 00:00:03.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO station (id, station_type_id, name, description, location_x, location_y, date_created, date_updated) VALUES
	('ddd3f222-590c-11eb-ae93-0242ac130002', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'Base Station One', 'Some description of Base Station One', 1, 1, '2021-01-01 00:00:01.000001+00', '2021-01-01 00:00:01.000001+00'),
    ('ee72a90c-590c-11eb-ae93-0242ac130002', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'Water Station One', 'Some description of Water Station One', 2, 2, '2021-01-01 00:00:02.000001+00', '2021-01-01 00:00:02.000001+00'),
    ('f676f266-590c-11eb-ae93-0242ac130002', '5c86bbaa-4ef8-11eb-ae93-0242ac130002', 'Plant Station One', 'Some description of Plant Station One', 3, 3, '2021-01-01 00:00:03.000001+00', '2021-01-01 00:00:03.000001+00'),
    ('feaa0806-590c-11eb-ae93-0242ac130002', '5c86bbaa-4ef8-11eb-ae93-0242ac130002', 'Plant Station Two', 'Some description of Plant Station Two', 4, 3, '2021-01-01 00:00:04.000001+00', '2021-01-01 00:00:04.000001+00'),
    ('0690d086-590d-11eb-ae93-0242ac130002', '5c86bbaa-4ef8-11eb-ae93-0242ac130002', 'Plant Station Three', 'Some description of Plant Station Three', 5, 3, '2021-01-01 00:00:05.000001+00', '2021-01-01 00:00:05.000001+00')
	ON CONFLICT DO NOTHING;

-- Create admin and regular Account with password "gophers"
INSERT INTO account (id, name, roles, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin', '{ADMIN,STATION}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2021-01-01 00:00:00', '2021-01-01 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Station', '{STATION}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2021-01-01 00:00:00', '2019-01-01 00:00:00')
	ON CONFLICT DO NOTHING;
`

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
