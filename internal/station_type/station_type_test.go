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

func TestStationType(t *testing.T) {
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

	update := station_type.UpdateStationType{
		Name: tests.StringPointer("Station Type 0"),
		Description: tests.StringPointer("Station type description 0"),
	}
	updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)

	if err := station_type.Update(ctx, db, st0.Id, update, updatedTime); err != nil {
		t.Fatalf("updating station type st0: %s", err)
	}

	saved, err := station_type.Retrieve(ctx, db, st0.Id)
	if err != nil {
		t.Fatalf("getting station type st0: %s", err)
	}

	// Check specified fields were updated. Make a copy of the original product
	// and change just the fields we expect then diff it with what was saved.
	want := *st0
	want.Name = "Station Type 0"
	want.Description = "Station type description 0"
	want.DateUpdated = updatedTime

	if diff := cmp.Diff(want, *saved); diff != "" {
		t.Fatalf("updated record did not match:\n%s", diff)
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
		t.Fatalf("expected station type list size %v, got %v", exp, got)
	}
}
