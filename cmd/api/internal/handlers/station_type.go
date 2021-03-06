package handlers

import (
	// Core packages
	"context"
	"log"
	"net/http"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/auth"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/web"
	"github.com/deezone/HydroBytes-BaseStation/internal/station_type"

	// Third party packages
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type StationType struct {
	db *sqlx.DB
	log *log.Logger
}

// Create decodes the body of a request to create a new station type. The full
// station type with generated fields is sent back in the response.
func (st *StationType) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.StationType.Create")
	defer span.End()

	var nst station_type.NewStationType

	if err := web.Decode(r, &nst); err != nil {
		return errors.Wrap(err, "decoding new station type")
	}

	stationType, err := station_type.Create(ctx, st.db, nst, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new station type")
	}

	return web.Respond(ctx, w, &stationType, http.StatusCreated)
}

// Delete removes a single station type identified by an ID in the request URL.
func (p *StationType) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.StationType.Delete")
	defer span.End()

	id := chi.URLParam(r, "id")

	if err := station_type.Delete(ctx, p.db, id); err != nil {
		switch err {
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting station type %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
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
func (st *StationType) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Product.List")
	defer span.End()

	list, err := station_type.List(ctx, st.db)
	if err != nil {
		return errors.Wrap(err, "getting station type list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Retrieve finds a station type identified by a station type ID in the request URL.
func (st *StationType) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.StationType.Retrieve")
	defer span.End()

	id := chi.URLParam(r, "id")

	stationTypes, err := station_type.Get(ctx, st.db, id)
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

	return web.Respond(ctx, w, stationTypes, http.StatusOK)
}

// Update decodes the body of a request to update an existing station type. The ID
// of the station type is part of the request URL.
func (st *StationType) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.StationType.Update")
	defer span.End()

	id := chi.URLParam(r, "id")

	var update station_type.UpdateStationType
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding station type update")
	}

	if err := station_type.Update(ctx, st.db, id, update, time.Now()); err != nil {
		switch err {
		case station_type.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating station type %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

/**
 * Station handlers
 */

// AddStation creates a new Station for a particular station_type. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (st *StationType) AddStation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Station.AddStation")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var ns station_type.NewStation
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new station")
	}

	stationTypeId := chi.URLParam(r, "id")

	station, err := station_type.AddStation(ctx, st.db, claims, ns, stationTypeId, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(ctx, w, station, http.StatusCreated)
}

// AdjustStation decodes the body of a request to update an existing station. The ID
// of the station is part of the request URL.
func (st *StationType) AdjustStation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Station.AdjustStation")
	defer span.End()

	id := chi.URLParam(r, "id")

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var update station_type.UpdateStation
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding station update")
	}

	if err := station_type.AdjustStation(ctx, st.db, claims, id, update, time.Now()); err != nil {
		switch err {
		case station_type.ErrStationNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case station_type.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating station %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// DeleteStation removes a single station identified by an ID in the request URL.
func (p *StationType) DeleteStation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Station.DeleteStation")
	defer span.End()

	id := chi.URLParam(r, "id")

	if err := station_type.DeleteStation(ctx, p.db, id); err != nil {
		switch err {
		case station_type.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting station %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// ListStations gets all sales for a particular station type.
func (st *StationType) ListStations(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Station.ListStations")
	defer span.End()

	id := chi.URLParam(r, "id")

	list, err := station_type.ListStations(ctx, st.db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// RetrieveStation finds a single station identified by an ID in the request URL.
func (st *StationType) RetrieveStation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Station.RetrieveStation")
	defer span.End()

	id := chi.URLParam(r, "id")

	station, err := station_type.GetStation(ctx, st.db, id)
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

	return web.Respond(ctx, w, station, http.StatusOK)
}
