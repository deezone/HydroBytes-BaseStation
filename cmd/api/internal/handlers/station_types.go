package handlers

import (
	// Core packages
	"encoding/json"
	"log"
	"net/http"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/station_types"

	// Third party packages
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type StationTypes struct {
	db *sqlx.DB
	log *log.Logger
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
func (st *StationTypes) List(w http.ResponseWriter, r *http.Request) {

	list, err := station_types.List(st.db)
	if err != nil {
		st.log.Println("listing station types", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// https://golang.org/pkg/encoding/json/#Marshal
	data, err := json.Marshal(list)
	if err != nil {
		st.log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// https://golang.org/pkg/net/http/#Request.Write
	if _, err := w.Write(data); err != nil {
		st.log.Println("error writing result", err)
	}
}

// Retrieve finds a single station type identified by an ID in the request URL.
func (st *StationTypes) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	station, err := station_types.Retrieve(st.db, id)
	if err != nil {
		st.log.Println("getting station type", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(station)
	if err != nil {
		st.log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		st.log.Println("error writing result", err)
	}
}
