package station_type_test

import (
	"context"
	// Core packages
	"testing"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/station_type"
	"github.com/deezone/HydroBytes-BaseStation/internal/schema"
	"github.com/deezone/HydroBytes-BaseStation/internal/tests"

	// Third-party packages
	"github.com/google/go-cmp/cmp"
)

func TestStationTypeCreateRetrieve(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	newP := station_type.NewStationType{
		Name:        "Base",
		Description: "Coordinator for all station types - monitor, command and control. Access point to public Intenet.",
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	st0, err := station_type.Create(ctx, db, newP, now)
	if err != nil {
		t.Fatalf("creating station type st0: %s", err)
	}

	st1, err := station_type.Retrieve(ctx, db, st0.Id)
	if err != nil {
		t.Fatalf("getting station type p0: %s", err)
	}

	if diff := cmp.Diff(st1, st0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}
}

func TestStationTypeList(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	sts, err := station_type.List(ctx, db)
	if err != nil {
		t.Fatalf("listing station types: %s", err)
	}
	if exp, got := 3, len(sts); exp != got {
		t.Fatalf("expected station types list size %v, got %v", exp, got)
	}
}
