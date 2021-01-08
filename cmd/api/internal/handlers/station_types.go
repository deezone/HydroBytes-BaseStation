package handlers

import (
	// Core packages
	"encoding/json"
	"log"
	"net/http"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/station_types"

	// Third party packages
	"github.com/jmoiron/sqlx"
)

type StationTypes struct {
	DB *sqlx.DB
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

	list, err := station_types.List(st.DB)
	if err != nil {
		log.Printf("error: listing station types: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// https://golang.org/pkg/encoding/json/#Marshal
	data, err := json.Marshal(list)
	if err != nil {
		log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// https://golang.org/pkg/net/http/#Request.Write
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}
}
