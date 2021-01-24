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

	st := StationType{db: db, log: log}

	app.Handle(http.MethodGet,    "/v1/station-types", st.List)
	app.Handle(http.MethodGet,    "/v1/station-type/{id}", st.Retrieve)
	app.Handle(http.MethodPost,   "/v1/station-type", st.Create)
	app.Handle(http.MethodPut,    "/v1/station-type/{id}", st.Update)
	app.Handle(http.MethodDelete, "/v1/station-type/{id}", st.Delete)

	app.Handle(http.MethodGet,    "/v1/station-type/{id}/stations", st.ListStations)
	app.Handle(http.MethodGet,    "/v1/station/{id}", st.RetrieveStation)
	app.Handle(http.MethodPost,   "/v1/station-type/{id}/station", st.AddStation)
	app.Handle(http.MethodPut,    "/v1/station/{id}", st.AdjustStation)
	app.Handle(http.MethodDelete, "/v1/station/{id}", st.DeleteStation)

	return app
}
