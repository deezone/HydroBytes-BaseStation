package handlers

import (
	// Core packages
	"log"
	"net/http"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/web"
	"github.com/deezone/HydroBytes-BaseStation/internal/station_type"

	// Third party packages
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type StationType struct {
	db *sqlx.DB
	log *log.Logger
}

// Create decodes the body of a request to create a new station type. The full
// station type with generated fields is sent back in the response.
func (st *StationType) Create(w http.ResponseWriter, r *http.Request) error {

	var nst station_type.NewStationType

	if err := web.Decode(r, &nst); err != nil {
		return errors.Wrap(err, "decoding new station type")
	}

	stationType, err := station_type.Create(r.Context(), st.db, nst, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new station type")
	}

	return web.Respond(r.Context(), w, &stationType, http.StatusCreated)
}

// Delete removes a single station type identified by an ID in the request URL.
func (p *StationType) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := station_type.Delete(r.Context(), p.db, id); err != nil {
		switch err {
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting station type %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
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
func (st *StationType) List(w http.ResponseWriter, r *http.Request) error {

	list, err := station_type.List(r.Context(), st.db)
	if err != nil {
		return errors.Wrap(err, "getting station type list")
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve finds all stations of a station type identified by a station type ID in the request URL.
func (st *StationType) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	stationTypes, err := station_type.Retrieve(r.Context(), st.db, id)
	if err != nil {
		switch err {
		case station_type.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting station type %q", id)
		}
	}

	return web.Respond(r.Context(), w, stationTypes, http.StatusOK)
}

// Update decodes the body of a request to update an existing station type. The ID
// of the station type is part of the request URL.
func (st *StationType) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update station_type.UpdateStationType
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding station type update")
	}

	if err := station_type.Update(r.Context(), st.db, id, update, time.Now()); err != nil {
		switch err {
		case station_type.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating station type %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

/**
 * Station handlers
 */

// AddStation creates a new Station for a particular station_type. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (st *StationType) AddStation(w http.ResponseWriter, r *http.Request) error {
	var ns station_type.NewStation
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new station")
	}

	stationTypeId := chi.URLParam(r, "id")

	station, err := station_type.AddStation(r.Context(), st.db, ns, stationTypeId, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(r.Context(), w, station, http.StatusCreated)
}

// Delete removes a single station identified by an ID in the request URL.
func (p *StationType) DeleteStation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := station_type.DeleteStation(r.Context(), p.db, id); err != nil {
		switch err {
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting station %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// ListStations gets all sales for a particular station type.
func (st *StationType) ListStations(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := station_type.ListStations(r.Context(), st.db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve finds a single station identified by an ID in the request URL.
func (st *StationType) RetrieveStation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	station, err := station_type.RetrieveStation(r.Context(), st.db, id)
	if err != nil {
		switch err {
		case station_type.ErrStationNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting station %q", id)
		}
	}

	return web.Respond(r.Context(), w, station, http.StatusOK)
}

// Update decodes the body of a request to update an existing station. The ID
// of the station is part of the request URL.
func (st *StationType) AdjustStation(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update station_type.UpdateStation
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding station update")
	}

	if err := station_type.AdjustStation(r.Context(), st.db, id, update, time.Now()); err != nil {
		switch err {
		case station_type.ErrStationNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating station %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}
