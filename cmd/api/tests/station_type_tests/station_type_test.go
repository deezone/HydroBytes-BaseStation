package station_type_tests

import (
	// Core Packages
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// NOTE: Models should not be imported, we want to test the exact JSON. We
	// make the comparison process easier using the go-cmp library.
	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/cmd/api/internal/handlers"
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
	test := tests.New(t)
	defer test.Teardown()

	stationTypeTests := StationTypeTests{app: handlers.API(test.Db, test.Log, test.Authenticator)}

	t.Run("List", stationTypeTests.List)
	t.Run("CreateRequiresFields", stationTypeTests.CreateRequiresFields)
	t.Run("StationTypeCRUD", stationTypeTests.StationTypeCRUD)
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
			"id":           "5c86bbaa-4ef8-11eb-ae93-0242ac130002",
			"name":         "Plant",
			"description":  "Monitors and reports plant health.",
			"stations":     float64(3),
			"date_created": "2021-01-01T00:00:03.000001Z",
			"date_updated": "2021-01-01T00:00:03.000001Z",
		},
		{
			"id":           "a2b0639f-2cc6-44b8-b97b-15d69dbb511e",
			"name":         "Base",
			"description":  "Coordinator for all station types - monitor, command and control. Access point to public Intenet.",
			"stations":     float64(1),
			"date_created": "2021-01-01T00:00:01.000001Z",
			"date_updated": "2021-01-01T00:00:01.000001Z",
		},
		{
			"id":           "72f8b983-3eb4-48db-9ed0-e45cc6bd716b",
			"name":         "Water",
			"description":  "Management of water resources. Controls water levels in resavour and impliments irrigation.",
			"stations":     float64(1),
			"date_created": "2021-01-01T00:00:02.000001Z",
			"date_updated": "2021-01-01T00:00:02.000001Z",
		},
	}

	if diff := cmp.Diff(expected, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

func (p *StationTypeTests) CreateRequiresFields(t *testing.T) {
	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("POST", "/v1/station-type", body)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusBadRequest, resp.Code)
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
			"stations":     float64(0),
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

	{ // UPDATE
		body := strings.NewReader(`{"name":"UPDATED stationtype0","description":"UPDATED Test description 0"}`)
		url := fmt.Sprintf("/v1/station-type/%s", actual["id"])
		req := httptest.NewRequest("PUT", url, body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		st.app.ServeHTTP(resp, req)

		if http.StatusNoContent != resp.Code {
			t.Fatalf("updating: expected status code %v, got %v", http.StatusNoContent, resp.Code)
		}

		// Retrieve updated record to be sure it worked.
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp = httptest.NewRecorder()

		st.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var updated map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		want := map[string]interface{}{
			"id":           actual["id"],
			"date_created": actual["date_created"],
			"date_updated": updated["date_updated"],
			"name":         "UPDATED stationtype0",
			"description":  "UPDATED Test description 0",
			"stations":     float64(0),
		}

		// Updated station type should match the one we created.
		if diff := cmp.Diff(want, updated); diff != "" {
			t.Fatalf("Retrieved station type should match created. Diff:\n%s", diff)
		}
	}

	{ // DELETE
		url := fmt.Sprintf("/v1/station-type/%s", actual["id"])
		req := httptest.NewRequest("DELETE", url, nil)
		resp := httptest.NewRecorder()

		st.app.ServeHTTP(resp, req)

		if http.StatusNoContent != resp.Code {
			t.Fatalf("updating: expected status code %v, got %v", http.StatusNoContent, resp.Code)
		}

		// Retrieve updated record to be sure it worked.
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp = httptest.NewRecorder()

		st.app.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusNotFound, resp.Code)
		}
	}
}
