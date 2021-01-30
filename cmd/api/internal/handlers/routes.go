package handlers

import (
	// Core packages
	"log"
	"net/http"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/mid"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/auth"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/web"

	// Third-party packages
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, log *log.Logger, authenticator *auth.Authenticator) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(log, mid.Logger(log), mid.Errors(log), mid.Metrics())

	{
		// Register health check handler. This route is not authenticated.
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		// Register account handlers.
		a := Account{db: db, authenticator: authenticator}
		app.Handle(http.MethodGet, "/v1/account/token", a.Token)
	}

	{
		// Register StationType handlers. Ensure all routes are authenticated.
		st := StationType{db: db, log: log}

		// StationType
		app.Handle(http.MethodGet,    "/v1/station-types",     st.List,     mid.Authenticate(authenticator))
		app.Handle(http.MethodGet,    "/v1/station-type/{id}", st.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost,   "/v1/station-type",      st.Create,
			mid.Authenticate(authenticator),
			mid.HasRole(auth.RoleAdmin),
		)
		app.Handle(http.MethodPut,    "/v1/station-type/{id}", st.Update,
			mid.Authenticate(authenticator),
			mid.HasRole(auth.RoleAdmin),
		)
		app.Handle(http.MethodDelete, "/v1/station-type/{id}", st.Delete,
			mid.Authenticate(authenticator),
			mid.HasRole(auth.RoleAdmin),
		)

		// Station
		app.Handle(http.MethodGet,    "/v1/station-type/{id}/stations", st.ListStations,    mid.Authenticate(authenticator))
		app.Handle(http.MethodGet,    "/v1/station/{id}",               st.RetrieveStation, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost,   "/v1/station-type/{id}/station",  st.AddStation,
			mid.Authenticate(authenticator),
			mid.HasRole(auth.RoleAdmin),
		)
		app.Handle(http.MethodPut,    "/v1/station/{id}",               st.AdjustStation,
			mid.Authenticate(authenticator),
			mid.HasRole(auth.RoleAdmin),
		)
		app.Handle(http.MethodDelete, "/v1/station/{id}",               st.DeleteStation,
			mid.Authenticate(authenticator),
			mid.HasRole(auth.RoleAdmin),
		)
	}

	return app
}
