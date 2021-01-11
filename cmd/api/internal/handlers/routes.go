package handlers

import (
	// Core packages
	"log"
	"net/http"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/web"

	// Third-party packages
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, log *log.Logger) http.Handler {

	app := web.NewApp(log)

	st := StationTypes{db: db, log: log}

	app.Handle(http.MethodGet, "/v1/station-types", st.List)
	app.Handle(http.MethodGet, "/v1/station-types/{id}", st.Retrieve)

	return app
}