package tests

import (
	// Core Packages
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	// NOTE: Models should not be imported, we want to test the exact JSON. We
	// make the comparison process easier using the go-cmp library.
	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/cmd/api/internal/handlers"
	"github.com/deezone/HydroBytes-BaseStation/internal/schema"
	"github.com/deezone/HydroBytes-BaseStation/internal/tests"

	// Third-party packages
	"github.com/google/go-cmp/cmp"
)

// TestStation runs a series of tests to exercise Station behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are run matters.
// One test may break if other tests are not run before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestStation(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	tests := StationTests{app: handlers.API(db, log)}

	t.Run("ListStations", tests.ListStations)
	t.Run("CreateRequiresFields", tests.CreateRequiresFields)
	t.Run("StationCRUD", tests.StationCRUD)
}

// StationTests holds methods for each station subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type StationTests struct {
	app http.Handler
}

func (st *StationTests) ListStations(t *testing.T) {

	// Get list of stations by the Plant StationType (5c86bbaa-4ef8-11eb-ae93-0242ac130002) as defined in the seed data
	req := httptest.NewRequest("GET", "/v1/station-type/5c86bbaa-4ef8-11eb-ae93-0242ac130002/stations", nil)
	resp := httptest.NewRecorder()

	st.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	expected := []map[string]interface{}{
		{
			"id":              "f676f266-590c-11eb-ae93-0242ac130002",
			"station_type_id": "5c86bbaa-4ef8-11eb-ae93-0242ac130002",
			"name":            "Plant Station One",
			"description":     "Some description of Plant Station One",
			"location_x":      float64(3),
			"location_y":      float64(3),
			"date_created":    "2021-01-01T00:00:03.000001Z",
			"date_updated":    "2021-01-01T00:00:03.000001Z",
		},
		{
			"id":              "feaa0806-590c-11eb-ae93-0242ac130002",
			"station_type_id": "5c86bbaa-4ef8-11eb-ae93-0242ac130002",
			"name":            "Plant Station Two",
			"description":     "Some description of Plant Station Two",
			"location_x":      float64(4),
			"location_y":      float64(3),
			"date_created":    "2021-01-01T00:00:04.000001Z",
			"date_updated":    "2021-01-01T00:00:04.000001Z",
		},
		{
			"id":              "0690d086-590d-11eb-ae93-0242ac130002",
			"station_type_id": "5c86bbaa-4ef8-11eb-ae93-0242ac130002",
			"name":            "Plant Station Three",
			"description":     "Some description of Plant Station Three",
			"location_x":      float64(5),
			"location_y":      float64(3),
			"date_created":    "2021-01-01T00:00:05.000001Z",
			"date_updated":    "2021-01-01T00:00:05.000001Z",
		},
	}

	if diff := cmp.Diff(expected, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

func (p *StationTests) CreateRequiresFields(t *testing.T) {
	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("POST", "/v1/station-type/5c86bbaa-4ef8-11eb-ae93-0242ac130002/station", body)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusBadRequest, resp.Code)
	}
}

func (st *StationTests) StationCRUD(t *testing.T) {
	var actual map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"name":"station0","description":"Test description 0", "location_x":123, "location_y": 789}`)

		// Create new station of type Plant StationType (5c86bbaa-4ef8-11eb-ae93-0242ac130002) as defined in the seed data
		req := httptest.NewRequest("POST", "/v1/station-type/5c86bbaa-4ef8-11eb-ae93-0242ac130002/station", body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		st.app.ServeHTTP(resp, req)

		if http.StatusCreated != resp.Code {
			t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if actual["id"] == "" || actual["id"] == nil {
			t.Fatal("expected non-empty station_type id")
		}
		if actual["date_created"] == "" || actual["date_created"] == nil {
			t.Fatal("expected non-empty station_type date_created")
		}
		if actual["date_updated"] == "" || actual["date_updated"] == nil {
			t.Fatal("expected non-empty station_type date_updated")
		}

		expected := map[string]interface{}{
			"id":           actual["id"],
			"date_created": actual["date_created"],
			"date_updated": actual["date_updated"],
			"name":         "station0",
			"description":  "Test description 0",
			"location_x":   float64(123),
			"location_y":   float64(789),
		}

		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Fatalf("Response did not match expected. Diff:\n%s", diff)
		}
	}
}
