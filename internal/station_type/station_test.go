package station_type_test

import (
	// Core packages
	"context"
	"testing"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/station_type"
	"github.com/deezone/HydroBytes-BaseStation/internal/tests"
)

func TestStation(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	ctx := context.Background()

	// Create two station types to work with.
	newStationTypeOne := station_type.NewStationType{
		Name:     "Station Type Name One",
		Description: "Station type one description test",
	}

	stationTypeOne, err := station_type.Create(ctx, db, newStationTypeOne, now)
	if err != nil {
		t.Fatalf("creating station type one: %s", err)
	}

	newStationTypeTwo := station_type.NewStationType{
		Name:     "Station Type Name Two",
		Description: "Station type two description test",

	}
	stationTypeTwo, err := station_type.Create(ctx, db, newStationTypeTwo, now)
	if err != nil {
		t.Fatalf("creating station type two: %s", err)
	}

	{ // Add and list

		ns := station_type.NewStation{
			Name: "",
			Description: "",
			LocationX: 7,
			LocationY: 6,
		}

		s, err := station_type.AddStation(ctx, db, ns, stationTypeOne.Id, now)
		if err != nil {
			t.Fatalf("adding test station one: %s", err)
		}

		// StationTypeOne should show the 1 station.
		stations, err := station_type.ListStations(ctx, db, stationTypeOne.Id)
		if err != nil {
			t.Fatalf("listing stations: %s", err)
		}
		if exp, got := 1, len(stations); exp != got {
			t.Fatalf("expected station list size %v, got %v", exp, got)
		}

		if exp, got := s.Id, stations[0].Id; exp != got {
			t.Fatalf("expected first station Id %v, got %v", exp, got)
		}

		// StationTypeTwo should have 0 stations.
		stations, err = station_type.ListStations(ctx, db, stationTypeTwo.Id)
		if err != nil {
			t.Fatalf("listing stations: %s", err)
		}
		if exp, got := 0, len(stations); exp != got {
			t.Fatalf("expected station list size %v, got %v", exp, got)
		}
	}
}
