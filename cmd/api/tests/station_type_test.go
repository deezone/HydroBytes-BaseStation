package tests

import (
	// Core Packages
	"encoding/json"
	"fmt"
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

// TestStationType runs a series of tests to exercise StationTypes behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are run matters.
// One test may break if other tests are not run before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestStationType(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	tests := StationTypeTests{app: handlers.API(db, log)}

	t.Run("List", tests.List)
	t.Run("StationTypeCRUD", tests.StationTypeCRUD)
}

// StationTypesTests holds methods for each station types subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type StationTypeTests struct {
	app http.Handler
}

func (st *StationTypeTests) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/station-types", nil)
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
			"id":           "a2b0639f-2cc6-44b8-b97b-15d69dbb511e",
			"name":         "Base",
			"description":  "Coordinator for all station types - monitor, command and control. Access point to public Intenet.",
			"date_created": "2021-01-01T00:00:01.000001Z",
			"date_updated": "2021-01-01T00:00:01.000001Z",
		},
		{
			"id":           "72f8b983-3eb4-48db-9ed0-e45cc6bd716b",
			"name":         "Water",
			"description":  "Management of water resources. Controls water levels in resavour and impliments irrigation.",
			"date_created": "2021-01-01T00:00:02.000001Z",
			"date_updated": "2021-01-01T00:00:02.000001Z",
		},
		{
			"id":           "5c86bbaa-4ef8-11eb-ae93-0242ac130002",
			"name":         "Plant",
			"description":  "Monitors and reports plant health.",
			"date_created": "2021-01-01T00:00:03.000001Z",
			"date_updated": "2021-01-01T00:00:03.000001Z",
		},
	}

	if diff := cmp.Diff(expected, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

func (st *StationTypeTests) StationTypeCRUD(t *testing.T) {
	var actual map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"name":"stationtype0","description":"Test description 0"}`)

		req := httptest.NewRequest("POST", "/v1/station-type", body)
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
			"name":         "stationtype0",
			"description":  "Test description 0",
		}

		if diff := cmp.Diff(expected, actual); diff != "" {
			t.Fatalf("Response did not match expected. Diff:\n%s", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/v1/station-type/%s", actual["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		st.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		// Fetched station type should match the one created.
		if diff := cmp.Diff(actual, fetched); diff != "" {
			t.Fatalf("Retrieved station type should match created. Diff:\n%s", diff)
		}
	}
}
