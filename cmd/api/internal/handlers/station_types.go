package handlers

import (
	"github.com/pkg/errors"
	// Core packages
	"log"
	"net/http"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/web"
	"github.com/deezone/HydroBytes-BaseStation/internal/station_types"

	// Third party packages
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type StationTypes struct {
	db *sqlx.DB
	log *log.Logger
}

// Create decodes the body of a request to create a new station type. The full
// station type with generated fields is sent back in the response.
func (st *StationTypes) Create(w http.ResponseWriter, r *http.Request) error {

	var nst station_types.NewStationTypes

	if err := web.Decode(r, &nst); err != nil {
		return errors.Wrap(err, "decoding new station tyoe")
	}

	station_type, err := station_types.Create(st.db, nst, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new station type")
	}

	return web.Respond(w, &station_type, http.StatusCreated)
}

/**
 * List StationTypes is a basic HTTP Handler that lists all of the station types in the HydroByte Automated Garden system.
 * A Handler is also know as a "controller" in the MVC pattern.
 *
 * The parameters must follow the signature defined in the http.HandlerFunc() adapter
 * (https://golang.org/src/net/http/server.go?s=97511:97566#L2034) to convert this method into a HTTP handler type.
 *
 * Note: If you open localhost:8000 in your browser, you may notice double requests being made. This happens because
 * the browser sends a request in the background for a website favicon. More the reason to use Postman to test!
 */
func (st *StationTypes) List(w http.ResponseWriter, r *http.Request) error {

	list, err := station_types.List(st.db)
	if err != nil {
		return errors.Wrap(err, "getting station tyoe list")
	}

	return web.Respond(w, list, http.StatusOK)
}

// Retrieve finds a single station type identified by an ID in the request URL.
func (st *StationTypes) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	station_type, err := station_types.Retrieve(st.db, id)
	if err != nil {
		switch err {
		case station_types.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case station_types.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting station type %q", id)
		}
	}

	return web.Respond(w, station_type, http.StatusOK)
}
