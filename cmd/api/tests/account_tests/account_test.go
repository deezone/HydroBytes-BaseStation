package account_tests

import (
	// Core packages
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/cmd/api/internal/handlers"
	"github.com/deezone/HydroBytes-BaseStation/internal/tests"
)

// TestAccount runs a series of tests to exercise Account behavior.
func TestAccount(t *testing.T) {
	test := tests.New(t)
	defer test.Teardown()

	ut := AccountTests{
		app:        handlers.API(test.Db, test.Log, test.Authenticator),
		adminToken: test.Token("Admin", "gophers"),
	}

	t.Run("TokenRequireAuth", ut.TokenRequireAuth)
	t.Run("TokenDenyUnknown", ut.TokenDenyUnknown)
	t.Run("TokenDenyBadPassword", ut.TokenDenyBadPassword)
	t.Run("TokenSuccess", ut.TokenSuccess)
}

// AccountTests holds methods for each account subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type AccountTests struct {
	app http.Handler
	adminToken string
}

// TokenRequireAuth ensures that requests with no authentication are denied.
func (at *AccountTests) TokenRequireAuth(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/account/token", nil)
	resp := httptest.NewRecorder()

	at.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusUnauthorized, resp.Code)
	}
}

// TokenDenyUnknown ensures that account with an unrecognized name aren't given a token.
func (at *AccountTests) TokenDenyUnknown(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/account/token", nil)
	resp := httptest.NewRecorder()

	req.SetBasicAuth("BadDude", "gophers")

	at.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusUnauthorized, resp.Code)
	}
}

// TokenDenyBadPassword ensures that a known account with a bad password is not authenticated.
func (at *AccountTests) TokenDenyBadPassword(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/account/token", nil)
	resp := httptest.NewRecorder()

	req.SetBasicAuth("Admin", "GOPHERS")

	at.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusUnauthorized, resp.Code)
	}
}

// TokenSuccess tests that a known account with a good password gets a token.
func (at *AccountTests) TokenSuccess(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/account/token", nil)
	resp := httptest.NewRecorder()

	req.SetBasicAuth("Admin", "gophers")

	at.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var got map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	if len(got) != 1 {
		t.Error("unexpected values in token response")
	}

	if got["token"] == "" {
		t.Fatal("token was not in response")
	}
}
